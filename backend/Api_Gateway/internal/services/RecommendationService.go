package services

import (
	"API_Service/internal/domains/models"
	errs "API_Service/internal/errors"
	"context"
	"fmt"
	"log/slog"

	rsv1 "github.com/dima11223432/recommendationService_protos/api/gen/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RecommendationService struct {
	RecommendationClient rsv1.RecommendationServiceClient
	log                  *slog.Logger
	authinterceptor      AuthInterceptor
}

func NewRecommendationService(recommendationClient rsv1.RecommendationServiceClient, log *slog.Logger, authinterceptor AuthInterceptor) *RecommendationService {
	return &RecommendationService{
		RecommendationClient: recommendationClient,
		log:                  log,
		authinterceptor:      authinterceptor,
	}
}

func (r *RecommendationService) GetAllParsingChannels(
	ctx context.Context,
) ([]string, error) {
	const op = "Api_Service.internal.services.RecommendationService.GetAllParsingChannels"

	parsingChannels, err := r.RecommendationClient.GetAllParsingChannels(
		ctx,
		&rsv1.GetAllParsingChannelsRequest{},
	)
	if err != nil {
		r.log.Error("failed to get all parsing channels", slog.String("op", op), slog.Any("err", err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return parsingChannels.Channels, nil
}

func (r *RecommendationService) AddNewUserCustomParsingChannels(ctx context.Context, channel string) error {
	const op = "Api_Service.internal.services.RecommendationService.AddUserCustomParsingChannels"
	userID, err := r.authinterceptor.GetUserIdFromContext(ctx)
	if err != nil {
		r.log.Error("invalid user id in context", slog.String("op", op), slog.Any("err", err))
		return fmt.Errorf("%s, %w", op, err)
	}
	if channel == "" {
		return fmt.Errorf("%s, %w", op, errs.ErrIsEmpty)
	}

	_, err = r.RecommendationClient.AddNewUserCustomParsingChannel(ctx, &rsv1.AddNewUserCustomParsingChannelRequest{
		ChannelUsername: channel,
		UserId:          userID,
	})
	if err != nil {
		r.log.Error("failed to add new user custom parsing channel", slog.String("op", op), slog.Any("err", err))
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (r *RecommendationService) AddNewParsingChannel(ctx context.Context, channel string, category string) error {
	if channel == "" {
		r.log.Error("channelUsername is empty")
		return fmt.Errorf("channel is empty")
	}
	const op = "Api_Service.internal.services.RecommendationService.AddNewParsingChannel"
	_, err := r.RecommendationClient.AddNewParsingChannel(ctx, &rsv1.AddNewParsingChannelRequest{
		ChannelUsername: channel,
		Category:        category,
	})
	if err != nil {
		st, ok := status.FromError(err)

		if ok && st.Code() == codes.AlreadyExists {
			return fmt.Errorf("%s, %w", op, errs.ErrChannelExists)
		}

		r.log.Error("failed to add new parsing channel", slog.String("op", op), slog.Any("err", err))
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (r *RecommendationService) DeleteUserCustomParsingChannel(
	ctx context.Context,
	channel string,
) error {
	const op = "Api_Service.internal.services.RecommendationService.DeleteUserCustomParsingChannel"
	userID, err := r.authinterceptor.GetUserIdFromContext(ctx)
	if err != nil {
		r.log.Error("Invalid user id in context",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = r.RecommendationClient.DeleteUserCustomParsingChannel(ctx,
		&rsv1.DeleteUserCustomParsingChannelRequest{
			UserId:          userID,
			ChannelUsername: channel,
		},
	)
	if err != nil {
		r.log.Error("Failed to delete user custom parsing channel",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *RecommendationService) DeleteParsingChannel(ctx context.Context, channel string) error {
	if channel == "" {
		r.log.Error("channelUsername is empty")
		return fmt.Errorf("channel is empty")
	}
	const op = "Api_Service.internal.services.RecommendationService.DeleteParsingChannel"
	_, err := r.RecommendationClient.DeleteParsingChannel(ctx, &rsv1.DeleteParsingChannelRequest{ChannelUsername: channel})
	if err != nil {
		r.log.Error("failed to delete parsing channel", slog.String("op", op), slog.Any("err", err))
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

func (r *RecommendationService) GetUserRecommendatedPosts(ctx context.Context, cursor *models.Cursor) ([]models.Post, *models.Cursor, error) {
	const op = "Api_Service.internal.services.RecommendationService.GetUserPriorityPosts"
	userID, err := r.authinterceptor.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("%s, %w", op, err)
	}
	var newCursor *rsv1.Cursor
	if cursor != nil {
		newCursor = &rsv1.Cursor{
			Score: int64(cursor.Score),
			Id:    cursor.ID,
		}
	}
	res, err := r.RecommendationClient.GetRecommendatedPosts(
		ctx,
		&rsv1.GetRecommendatedPostsRequest{
			UserId: userID,
			Cursor: newCursor,
		},
	)
	var nextCursor *models.Cursor
	if res.NextCursor != nil {
		nextCursor = &models.Cursor{
			Score: float64(res.NextCursor.Score),
			ID:    res.NextCursor.Id,
		}
	}

	if err != nil {
		return nil, nil, fmt.Errorf("%s, %w", op, err)
	}
	var posts []models.Post
	for _, post := range res.Posts {
		var stocks []models.Stock
		for _, stock := range post.Stocks {
			stocks = append(stocks, models.Stock{
				StockName: stock.StockName,
				Side:      stock.Side,
			})
		}
		posts = append(posts, models.Post{
			Stocks:          stocks,
			PostText:        post.PostText,
			PostURI:         post.PostUri,
			Reasoning:       post.Reasoning,
			ChannelUsername: post.ChannelUsername,
			Date:            post.Date.AsTime(),
		})
	}
	return posts, nextCursor, nil

}

func (s *RecommendationService) GetUserPriorityChannels(ctx context.Context) ([]string, error) {
	const op = "Api_Service.internal.services.RecommendationService.GetUserPriorityChannels"

	userID, err := s.authinterceptor.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	res, err := s.RecommendationClient.GetUserPriorityChannels(
		ctx,
		&rsv1.GetUserPriorityChannelsRequest{
			UserId: userID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	channels := make([]string, 0, len(res.Channels))
	for _, channel := range res.GetChannels() {
		channels = append(channels, channel)
	}
	return channels, nil
}

func (r *RecommendationService) GetAllParsingChannelsWithCategories(ctx context.Context) (map[string][]string, error) {
	const op = "Api_Service.internal.services.RecommendationService.GetAllParsingChannelsWithCategories"
	res, err := r.RecommendationClient.GetAllParsingChannelsWithCategories(
		ctx,
		&rsv1.GetAllParsingChannelsWithCategoriesRequest{},
	)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	respChannels := res.GetChannels()
	filteredChannels := make(map[string][]string, 0)
	for category, channels := range respChannels {
		convertedChannels := make([]string, 0, len(channels.Usernames))
		for _, channel := range channels.Usernames {
			convertedChannels = append(convertedChannels, channel)
		}
		filteredChannels[category] = convertedChannels
	}
	return filteredChannels, nil
}

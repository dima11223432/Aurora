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

// RecommendationService handles recommendation and parsing channel operations.
// It communicates with the Recommendation Service via gRPC to fetch
// recommended posts and manage parsing channels (default and user-custom).
type RecommendationService struct {
	RecommendationClient rsv1.RecommendationServiceClient
	log                  *slog.Logger
	authinterceptor      AuthInterceptor
}

// NewRecommendationService creates a new RecommendationService instance.
func NewRecommendationService(recommendationClient rsv1.RecommendationServiceClient, log *slog.Logger, authinterceptor AuthInterceptor) *RecommendationService {
	return &RecommendationService{
		RecommendationClient: recommendationClient,
		log:                  log,
		authinterceptor:      authinterceptor,
	}
}

// GetAllDefaultParsingChannels returns the list of all system-wide default parsing channels.
func (r *RecommendationService) GetAllDefaultParsingChannels(
	ctx context.Context,
) ([]string, error) {
	const op = "Api_Service.internal.services.RecommendationService.GetAllDefaultParsingChannels"

	parsingChannels, err := r.RecommendationClient.GetAllDefaultParsingChannels(
		ctx,
		&rsv1.GetAllDefaultParsingChannelsRequest{},
	)
	if err != nil {
		r.log.Error("failed to get all default parsing channels", slog.String("op", op), slog.Any("err", err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return parsingChannels.Channels, nil
}

// AddNewUserCustomParsingChannels adds a custom parsing channel for the authenticated user.
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
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.AlreadyExists {
				return fmt.Errorf("%s:%w", op, errs.ErrChannelExists)
			}
		}
		r.log.Error("failed to add new user custom parsing channel", slog.String("op", op), slog.Any("err", err))
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

// AddNewDefaultParsingChannel adds a new system-wide default parsing channel with a category.
func (r *RecommendationService) AddNewDefaultParsingChannel(ctx context.Context, channel string, category string) error {
	if channel == "" {
		r.log.Error("channelUsername is empty")
		return fmt.Errorf("channel is empty")
	}
	const op = "Api_Service.internal.services.RecommendationService.AddNewDefaultParsingChannel"
	_, err := r.RecommendationClient.AddNewDefaultParsingChannel(ctx, &rsv1.AddNewDefaultParsingChannelRequest{
		ChannelUsername: channel,
		Category:        category,
	})
	if err != nil {
		st, ok := status.FromError(err)

		if ok && st.Code() == codes.AlreadyExists {
			return fmt.Errorf("%s, %w", op, errs.ErrChannelExists)
		}

		r.log.Error("failed to add new default parsing channel", slog.String("op", op), slog.Any("err", err))
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

// DeleteUserCustomParsingChannel removes a user's custom parsing channel.
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

// DeleteDefaultParsingChannel removes a system-wide default parsing channel.
func (r *RecommendationService) DeleteDefaultParsingChannel(ctx context.Context, channel string) error {
	if channel == "" {
		r.log.Error("channelUsername is empty")
		return fmt.Errorf("channel is empty")
	}
	const op = "Api_Service.internal.services.RecommendationService.DeleteDefaultParsingChannel"
	_, err := r.RecommendationClient.DeleteDefaultParsingChannel(ctx, &rsv1.DeleteDefaultParsingChannelRequest{
		ChannelUsername: channel,
	})
	if err != nil {
		r.log.Error("failed to delete default parsing channel", slog.String("op", op), slog.Any("err", err))
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

// GetUserRecommendatedPosts fetches recommended posts for the authenticated user
// with cursor-based pagination. Returns posts, next cursor, and error.
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
	if err != nil {
		return nil, nil, fmt.Errorf("%s, %w", op, err)
	}

	var nextCursor *models.Cursor
	if res.NextCursor != nil {
		nextCursor = &models.Cursor{
			Score: float64(res.NextCursor.Score),
			ID:    res.NextCursor.Id,
		}
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

// GetUserPriorityChannels returns the priority channels for the authenticated user.
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

// GetAllUserCustomParsingChannels returns all custom parsing channels for the authenticated user.
func (r *RecommendationService) GetAllUserCustomParsingChannels(ctx context.Context) ([]string, error) {
	const op = "Api_Service.internal.services.RecommendationService.GetAllUserCustomParsingChannel"
	userID, err := r.authinterceptor.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := r.RecommendationClient.GetAllUserCustomParsingChannels(
		ctx,
		&rsv1.GetAllUserCustomParsingChannelsRequest{
			UserId: userID,
		})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	channels := make([]string, 0, len(res.Channels))
	for _, channel := range res.GetChannels() {
		channels = append(channels, channel)
	}
	return channels, nil
}

// GetAllDefaultParsingChannelsWithCategories returns default parsing channels
// grouped by their categories.
func (r *RecommendationService) GetAllDefaultParsingChannelsWithCategories(ctx context.Context) (map[string][]string, error) {
	const op = "Api_Service.internal.services.RecommendationService.GetAllDefaultParsingChannelsWithCategories"
	res, err := r.RecommendationClient.GetAllDefaultParsingChannelsWithCategories(
		ctx,
		&rsv1.GetAllDefaultParsingChannelsWithCategoriesRequest{},
	)
	if err != nil {
		r.log.Error("failed to get all default parsing channels", slog.String("op", op), slog.Any("err", err))
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	respChannels := res.GetChannels()
	filteredChannels := make(map[string][]string)
	for category, channels := range respChannels {
		convertedChannels := make([]string, 0, len(channels.Usernames))
		for _, channel := range channels.Usernames {
			convertedChannels = append(convertedChannels, channel)
		}
		filteredChannels[category] = convertedChannels
	}
	return filteredChannels, nil
}

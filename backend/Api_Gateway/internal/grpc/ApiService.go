package grpc

import (
	v1 "API_Service/api/gen/v1"

	custom_errors "API_Service/internal/custom_errors"
	"API_Service/internal/domains/models"
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Auth interface {
	Login(ctx context.Context, telegram_id int64, username string, firstName string, lastName string, appId int64) (string, error)
	IsAdmin(ctx context.Context, telegram_id int64) (bool, error)
	SetPriorityChannels(ctx context.Context, channels []string) error
	DeletePriorityChannels(ctx context.Context, channels []string) error
}

type RecommendatinService interface {
	GetUserRecommendatedPosts(ctx context.Context, cursor *models.Cursor) ([]models.Post, *models.Cursor, error)
	GetAllParsingChannels(ctx context.Context) ([]string, error)
}

type ApiService struct {
	v1.UnimplementedApiServiceServer
	recommendatinService RecommendatinService
	auth                 Auth
}

func RegisterGrpcServer(gRPC *grpc.Server, auth Auth, recommendationService RecommendatinService) {
	v1.RegisterApiServiceServer(gRPC, &ApiService{
		auth:                 auth,
		recommendatinService: recommendationService,
	})
}

func (a *ApiService) Login(
	ctx context.Context,
	req *v1.LoginRequest,
) (*v1.LoginResponse, error) {

	token, err := a.auth.Login(
		ctx,
		req.GetTelegramId(),
		req.GetUsername(),
		req.GetFirstName(),
		req.GetLastName(),
		req.AppId,
	)
	if err != nil {
		logrus.WithError(err).Error("login failed")
		return nil, err
	}

	return &v1.LoginResponse{
		Token: token,
	}, nil
}

func (a *ApiService) DeletePriorityChannels(
	ctx context.Context,
	req *v1.DeletePriorityChannelRequest,
) (*v1.DeletePriorityChannelResponse, error) {
	err := a.auth.DeletePriorityChannels(ctx, req.GetChannels())
	if err != nil {
		status.Error(codes.Internal, err.Error())
		return nil, err
	}
	return nil, nil
}

func (a *ApiService) SetPriorityChannels(
	ctx context.Context,
	req *v1.SetPriorityChannelsRequest,
) (*v1.SetPriorityChannelsResponse, error) {
	channels := req.GetPriorityChannels()

	err := a.auth.SetPriorityChannels(ctx, channels)
	if err != nil {
		if errors.Is(err, custom_errors.ErrChannelExists) {
			return nil, status.Error(codes.AlreadyExists, "channel already exists")
		}
		return nil, status.Error(status.Code(err), err.Error())
	}

	return &v1.SetPriorityChannelsResponse{}, nil
}

func (a *ApiService) IsAdmin(
	ctx context.Context,
	req *v1.IsAdminRequest,
) (*v1.IsAdminResponse, error) {
	isAdmin, err := a.auth.IsAdmin(ctx, req.TelegramId)
	if err != nil {
		return &v1.IsAdminResponse{}, err
	}

	return &v1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (a *ApiService) GetRecommendatedPosts(
	ctx context.Context,
	req *v1.GetRecommendatedPostsRequest,
) (*v1.GetRecommendatedPostsResponse, error) {
	const op = "backend/Api_Gateway/internal/grpc/ApiService.go"
	var cursor *models.Cursor
	if req.GetCursor() != nil {

		cursor = &models.Cursor{
			Score: float64(req.GetCursor().Score),
			ID:    req.GetCursor().Id,
		}
	}
	posts, nextCursor, err := a.recommendatinService.GetUserRecommendatedPosts(ctx, cursor)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts: %v", err)
	}

	var NextCursor *v1.Cursor
	if nextCursor != nil {
		NextCursor = &v1.Cursor{
			Score: int64(nextCursor.Score),
			Id:    nextCursor.ID,
		}
	}

	var postList []*v1.Post
	for _, post := range posts {
		var stocks []*v1.Stock
		for _, stock := range post.Stocks {
			stocks = append(stocks, &v1.Stock{
				StockName: stock.StockName,
				Side:      stock.Side})
		}
		postList = append(postList, &v1.Post{
			Stocks:          stocks,
			PostText:        post.PostText,
			PostUri:         post.PostURI,
			ChannelUsername: post.ChannelUsername,
			Reasoning:       post.Reasoning,
			Date:            timestamppb.New(post.Date),
		})
	}
	return &v1.GetRecommendatedPostsResponse{
		Posts:      postList,
		NextCursor: NextCursor,
	}, nil
}

func (a *ApiService) GetAllParsingChannels(
	ctx context.Context,
	_ *v1.GetAllParsingChannelsRequest,
) (*v1.GetAllParsingChannelsResponse, error) {
	channels, err := a.recommendatinService.GetAllParsingChannels(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get parsing channels")
	}

	return &v1.GetAllParsingChannelsResponse{
		Channels: channels,
	}, nil
}

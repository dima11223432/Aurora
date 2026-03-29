package grpc

import (
	v1 "API_Service/api/gen/v1"

	custom_errors "API_Service/internal/custom_errors"
	"API_Service/internal/domains/models"
	"context"
	"errors"
	"fmt"

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
}

type RecommendatinService interface {
	GetUserRecommendatedPosts(ctx context.Context) ([]models.Post, *models.Cursor, error)
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

	if err != nil {
		logrus.WithError(err).Error("failed to parse jwt user id from context")
		return nil, fmt.Errorf("fail to parse jwt")
	}

	return &v1.LoginResponse{
		Token: token,
	}, nil
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
	posts, nextCursor, err := a.recommendatinService.GetUserRecommendatedPosts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts: %v", err)
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
			Reasoning:       "",
			Date:            timestamppb.New(post.Date),
		})
	}
	return &v1.GetRecommendatedPostsResponse{
		Posts: postList,
		NextCursor: &v1.Cursor{
			Score: int64(nextCursor.Score),
			Id:    nextCursor.ID,
		},
	}, nil
}

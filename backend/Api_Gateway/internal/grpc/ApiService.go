// Package grpc implements the gRPC server for the API Gateway service.
// It defines the ApiService struct that handles incoming gRPC requests,
// delegating business logic to the Auth and RecommendationService interfaces.
package grpc

import (
	v1 "API_Service/api/gen/v1"

	custom_errors "API_Service/internal/custom_errors"
	"API_Service/internal/domains/models"
	errs "API_Service/internal/errors"
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Auth defines the authentication operations available to the gRPC handler.
type Auth interface {
	Login(ctx context.Context, telegram_id int64, username string, firstName string, lastName string, appId int64) (string, error)
	IsAdmin(ctx context.Context, telegram_id int64) (bool, error)
	SetPriorityChannels(ctx context.Context, channels []string) error
	ConnectWallet(ctx context.Context, wallet string) error
	DeletePriorityChannels(ctx context.Context, channels []string) error
	IsAdminByContext(ctx context.Context) (bool, error)
}

// RecommendationService defines the recommendation operations available to the gRPC handler.
type RecommendationService interface {
	GetUserRecommendatedPosts(ctx context.Context, cursor *models.Cursor) ([]models.Post, *models.Cursor, error)
	GetAllDefaultParsingChannels(ctx context.Context) ([]string, error)
	DeleteDefaultParsingChannel(ctx context.Context, channel string) error
	DeleteUserCustomParsingChannel(ctx context.Context, channel string) error
	AddNewDefaultParsingChannel(ctx context.Context, channel string, category string) error
	GetAllUserCustomParsingChannels(ctx context.Context) ([]string, error)
	GetUserPriorityChannels(ctx context.Context) ([]string, error)
	AddNewUserCustomParsingChannels(ctx context.Context, channel string) error
	GetAllDefaultParsingChannelsWithCategories(ctx context.Context) (map[string][]string, error)
}

// ApiService implements the v1.ApiServiceServer gRPC interface.
// It delegates requests to the Auth and RecommendationService layers.
type ApiService struct {
	v1.UnimplementedApiServiceServer
	recommendationService RecommendationService
	auth                  Auth
}

// RegisterGrpcServer registers the ApiService with the given gRPC server.
func RegisterGrpcServer(gRPC *grpc.Server, auth Auth, recommendationService RecommendationService) {
	v1.RegisterApiServiceServer(gRPC, &ApiService{
		auth:                  auth,
		recommendationService: recommendationService,
	})
}

// Login handles user login requests via gRPC.
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

// DeletePriorityChannels handles priority channel deletion requests via gRPC.
func (a *ApiService) DeletePriorityChannels(
	ctx context.Context,
	req *v1.DeletePriorityChannelRequest,
) (*v1.DeletePriorityChannelResponse, error) {
	err := a.auth.DeletePriorityChannels(ctx, req.GetChannels())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &v1.DeletePriorityChannelResponse{}, nil
}

// SetPriorityChannels handles priority channel setting requests via gRPC.
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.SetPriorityChannelsResponse{}, nil
}

// IsAdmin handles admin status check requests via gRPC.
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

func (a *ApiService) IsAdminByContext(
	ctx context.Context,
	req *v1.IsAdminRequest,
) *v1.IsAdminByContextResponce {
}

// GetRecommendatedPosts handles recommended posts retrieval via gRPC.
func (a *ApiService) GetRecommendatedPosts(
	ctx context.Context,
	req *v1.GetRecommendatedPostsRequest,
) (*v1.GetRecommendatedPostsResponse, error) {
	var cursor *models.Cursor
	if req.GetCursor() != nil {
		cursor = &models.Cursor{
			Score: float64(req.GetCursor().Score),
			ID:    req.GetCursor().Id,
		}
	}
	posts, nextCursor, err := a.recommendationService.GetUserRecommendatedPosts(ctx, cursor)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts: %v", err)
	}

	var protoNextCursor *v1.Cursor
	if nextCursor != nil {
		protoNextCursor = &v1.Cursor{
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
		NextCursor: protoNextCursor,
	}, nil
}

// GetAllUserCustomParsingChannels handles retrieval of user custom parsing channels via gRPC.
func (a *ApiService) GetAllUserCustomParsingChannels(
	ctx context.Context,
	_ *v1.GetAllUserCustomParsingChannelsRequest,
) (*v1.GetAllUserCustomParsingChannelsResponse, error) {
	channels, err := a.recommendationService.GetAllUserCustomParsingChannels(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get parsing channels")
	}
	return &v1.GetAllUserCustomParsingChannelsResponse{
		Channels: channels,
	}, nil
}

// GetAllDefaultParsingChannels handles retrieval of all default parsing channels via gRPC.
func (a *ApiService) GetAllDefaultParsingChannels(
	ctx context.Context,
	_ *v1.GetAllDefaultParsingChannelsRequest,
) (*v1.GetAllDefaultParsingChannelsResponse, error) {
	channels, err := a.recommendationService.GetAllDefaultParsingChannels(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get parsing channels")
	}

	return &v1.GetAllDefaultParsingChannelsResponse{
		Channels: channels,
	}, nil
}

// AddNewUserCustomParsingChannel handles adding a new user custom parsing channel via gRPC.
func (a *ApiService) AddNewUserCustomParsingChannel(
	ctx context.Context,
	req *v1.AddNewUserCustomParsingChannelRequest,
) (*v1.AddNewUserCustomParsingChannelResponse, error) {
	channel := req.GetChannelUsername()
	err := a.recommendationService.AddNewUserCustomParsingChannels(ctx, channel)
	if err != nil {
		if errors.Is(err, errs.ErrChannelExists) {
			return nil, status.Error(codes.AlreadyExists, "channel already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &v1.AddNewUserCustomParsingChannelResponse{}, nil
}

// AddNewDefaultParsingChannel handles adding a new default parsing channel via gRPC.
// Requires admin privileges.
func (a *ApiService) AddNewDefaultParsingChannel(
	ctx context.Context,
	req *v1.AddNewDefaultParsingChannelRequest,
) (*v1.AddNewDefaultParsingChannelResponse, error) {
	isAdmin, err := a.auth.IsAdminByContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid JWT")
	}
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}
	err = a.recommendationService.AddNewDefaultParsingChannel(ctx, req.ChannelUsername, req.Category)
	if err != nil {
		if errors.Is(err, errs.ErrChannelExists) {
			return nil, status.Error(codes.AlreadyExists, "channel already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &v1.AddNewDefaultParsingChannelResponse{}, nil
}

// DeleteUserCustomParsingChannel handles deletion of a user custom parsing channel via gRPC.
func (a *ApiService) DeleteUserCustomParsingChannel(
	ctx context.Context,
	req *v1.DeleteUserCustomParsingChannelRequest,
) (*v1.DeleteUserCustomParsingChannelResponse, error) {
	err := a.recommendationService.DeleteUserCustomParsingChannel(ctx, req.ChannelUsername)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete channel")
	}
	return &v1.DeleteUserCustomParsingChannelResponse{}, nil
}

// DeleteDefaultParsingChannel handles deletion of a default parsing channel via gRPC.
// Requires admin privileges.
func (a *ApiService) DeleteDefaultParsingChannel(
	ctx context.Context,
	req *v1.DeleteDefaultParsingChannelRequest,
) (*v1.DeleteDefaultParsingChannelResponse, error) {

	isAdmin, err := a.auth.IsAdminByContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid JWT")
	}
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}
	err = a.recommendationService.DeleteDefaultParsingChannel(ctx, req.ChannelUsername)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete channel")
	}
	return &v1.DeleteDefaultParsingChannelResponse{}, nil
}

// GetUserPriorityChannels handles retrieval of user priority channels via gRPC.
func (a *ApiService) GetUserPriorityChannels(
	ctx context.Context,
	_ *v1.GetUserPriorityChannelsRequest,
) (*v1.GetUserPriorityChannelsResponse, error) {

	channels, err := a.recommendationService.GetUserPriorityChannels(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get priority channels")
	}
	return &v1.GetUserPriorityChannelsResponse{
		Channels: channels,
	}, nil
}

// GetAllDefaultParsingChannelsWithCategories handles retrieval of default parsing
// channels grouped by category via gRPC.
func (a *ApiService) GetAllDefaultParsingChannelsWithCategories(
	ctx context.Context,
	_ *v1.GetAllDefaultParsingChannelsWithCategoriesRequest,
) (*v1.GetAllDefaultParsingChannelsWithCategoriesResponse, error) {

	channelsWithCategories, err := a.recommendationService.GetAllDefaultParsingChannelsWithCategories(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get channels with categories")
	}

	protoChannels := make(map[string]*v1.ChannelList)
	for category, channels := range channelsWithCategories {
		protoChannels[category] = &v1.ChannelList{
			Usernames: channels,
		}
	}

	return &v1.GetAllDefaultParsingChannelsWithCategoriesResponse{
		Channels: protoChannels,
	}, nil
}

// ConnectWallet handles wallet connection requests via gRPC.
func (a *ApiService) ConnectWallet(
	ctx context.Context,
	req *v1.ConnectWalletRequest,
) (*v1.ConnectWalletResponse, error) {
	err := a.auth.ConnectWallet(ctx, req.WalletAddress)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &v1.ConnectWalletResponse{}, nil
}

package grpc

import (
	v1 "API_Service/api/gen/v1"
	"API_Service/internal/domains/models"
	services "API_Service/internal/services"
	"context"
	"fmt"

	"API_Service/internal/grpc/AuthInterceptor"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, telegram_id int64, username string, firstName string, lastName string, appId int64) (string, error)
	IsAdmin(ctx context.Context, telegram_id int64) (bool, error)
	SetPriorityChannels(ctx context.Context, channels []string) error
	DeletePriorityChannels(ctx context.Context, channels []string) error
	IsAdminByContext(ctx context.Context) (bool, error)
}

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

type ApiService struct {
	v1.UnimplementedApiServiceServer
	recommendationService RecommendationService
	auth                  Auth
}

func RegisterGrpcServer(gRPC *grpc.Server, auth Auth, recommendationService RecommendationService) {
	v1.RegisterApiServiceServer(gRPC, &ApiService{
		auth:                  auth,
		recommendationService: recommendationService,
	})
}

func (a *ApiService) Register(
	ctx context.Context,
	req *v1.RegisterRequest,
) (*v1.RegisterResponse, error) {
	userId, err := a.auth.Register(ctx, req.GetEmail(), req.GetPassword(), req.GetIsAdmin())

	if err != nil {
		return &v1.RegisterResponse{}, err
	}

	return &v1.RegisterResponse{
		UserId: userId,
	}, nil

}

func (a *ApiService) Login(
	ctx context.Context,
	req *v1.LoginRequest,
) (*v1.LoginResponse, error) {
	token, err := a.auth.Login(ctx, req.Email, req.Password, req.AppId)

	telegramID, err := authinterceptor.GetUserIdFromContext(ctx)
	if err != nil {
		return &v1.LoginResponse{}, err
	}
	return &v1.LoginResponse{
		Token: token,
	}, nil

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

func (a *ApiService) IsAdmin(
	ctx context.Context,
	req *v1.IsAdminRequest,
) (*v1.IsAdminResponse, error) {
	isAdmin, err := a.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		return &v1.IsAdminResponse{}, err
	}

	return &v1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (a *ApiService) UpdateGroup(
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
		return nil, err
	}
	return &v1.UpdateGroupResponse{
		Success: true,
	}, nil
}

func (a *ApiService) DeleteGroup(
	ctx context.Context,
	req *v1.DeleteGroupRequest,
) (*v1.DeleteGroupResponse, error) {

	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id must be > 0")
	}

	var protoNextCursor *v1.Cursor
	if nextCursor != nil {
		protoNextCursor = &v1.Cursor{
			Score: int64(nextCursor.Score),
			Id:    nextCursor.ID,
		}
	}

	return &v1.DeleteGroupResponse{
		Success: true,
	}, nil
}

func (a *ApiService) CreateGroup(
	ctx context.Context,
	req *v1.CreateGroupRequest,
) (*v1.CreateGroupResponse, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("it is no JWT in metadata")
	}

	fmt.Println("Metadata", md)

	auth := md.Get("authorization")
	fmt.Println("auth", auth)

	id, err := a.api.CreateGroup(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	return &v1.CreateGroupResponse{
		Id:   id,
		Name: req.Name,
	}, nil
}

func (a *ApiService) DeleteContact(
	ctx context.Context,
	req *v1.DeleteContactRequest,
) (*v1.DeleteContactResponse, error) {
	err := a.api.DeleteContact(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteContactResponse{
		Success: true,
	}, nil
}

func (a *ApiService) GetAllContactsByGroupID(
	ctx context.Context,
	req *v1.GetAllContactsByGroupIDRequest,
) (*v1.GetAllContactsByGroupIDResponse, error) {
	contacts, err := a.api.GetAllContactsByGroupID(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}
	pbContact := make([]*v1.Contact, 0, len(contacts))

	for _, contact := range contacts {
		pbContact = append(pbContact, &v1.Contact{
			Id:    contact.ID,
			Email: contact.Email,
		})
	}
	return &v1.GetRecommendatedPostsResponse{
		Posts:      postList,
		NextCursor: protoNextCursor,
	}, nil
}

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

func (a *ApiService) GetAllDefaultParsingChannels(
	ctx context.Context,
	_ *v1.GetAllDefaultParsingChannelsRequest,
) (*v1.GetAllDefaultParsingChannelsResponse, error) {
	channels, err := a.recommendationService.GetAllDefaultParsingChannels(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.GetAllDefaultParsingChannelsResponse{
		Channels: channels,
	}, nil
}

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

func (a *ApiService) GetUserPriorityChannels(
	ctx context.Context,
	_ *v1.GetUserPriorityChannelsRequest,
) (*v1.GetUserPriorityChannelsResponse, error) {

	channels, err := a.recommendationService.GetUserPriorityChannels(ctx)
	if err != nil {
		return nil, err
	}

	return &v1.CreateNotificationResponse{
		Id:    id,
		Title: req.Title,
		Text:  req.Text,
	}, nil
}

func (a *ApiService) SendNotification(ctx context.Context, req *v1.SendNotificationRequest) (*v1.SendNotificationResponse, error) {

	userId, uerr := authinterceptor.GetUserIdFromContext(ctx)
	if uerr != nil {
		return nil, status.Errorf(codes.Internal, "fail to get user id")
	}
	logrus.Info("UserId = ", userId)
	EventId, err := a.api.SendNotification(ctx, req.NotificationId)

	if err != nil {
		return nil, err
	}

	return &v1.SendNotificationResponse{
		EventId: EventId,
		Status:  "ok",
	}, nil
}

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

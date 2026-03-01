package grpc

import (
	v1 "API_Service/api/gen/v1"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Auth interface {
	Login(ctx context.Context, telegram_id int64, username string, firstName string, lastName string, appId int64) (string, error)
	IsAdmin(ctx context.Context, telegram_id int64) (bool, error)
	SetPriorityChannels(ctx context.Context, channels []string) (int32, error)
}

type ApiService struct {
	v1.UnimplementedApiServiceServer
	auth Auth
}

func RegisterGrpcServer(gRPC *grpc.Server, auth Auth) {
	v1.RegisterApiServiceServer(gRPC, &ApiService{
		auth: auth,
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

	logrus.WithFields(logrus.Fields{
		"telegram_id": req.TelegramId,
		"username":    req.GetUsername(),
		"app_id":      req.AppId,
	}).Info("user logged in successfully")

	return &v1.LoginResponse{
		Token: token,
	}, nil
}

func (a *ApiService) SetPriorityChannels(
	ctx context.Context,
	req *v1.SetPriorityChannelsRequest,
) (*v1.SetPriorityChannelsResponse, error) {

	return nil, nil
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

func (a *ApiService) UpdateGroup(
	ctx context.Context,
	req *v1.UpdateGroupRequest,
) (*v1.UpdateGroupResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id must be > 0")
	}

	err := a.api.UpdateGroup(ctx, req.Id, req.Name)
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

	err := a.api.DeleteGroup(ctx, req.Id)
	if err != nil {
		return &v1.DeleteGroupResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "fail to delete group")
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
	return &v1.GetAllContactsByGroupIDResponse{
		Contacts: pbContact,
	}, nil
}

func (a *ApiService) UpdateContact(ctx context.Context, req *v1.UpdateContactRequest) (*v1.UpdateContactResponse, error) {
	err := a.api.UpdateContact(ctx, req.Id, req.Email)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateContactResponse{
		Success: true,
	}, nil
}

func (a *ApiService) CreateContact(ctx context.Context, req *v1.CreateContactRequest) (*v1.CreateContactResponse, error) {
	id, err := a.api.CreateContact(ctx, req.Email, req.GroupId)
	if err != nil {
		return nil, err
	}

	return &v1.CreateContactResponse{
		Id: id,
	}, nil
}

func (a *ApiService) CreateNotification(ctx context.Context, req *v1.CreateNotificationRequest) (*v1.CreateNotificationResponse, error) {
	id, err := a.api.CreateNotification(ctx, req.Title, req.Text, req.GroupId)
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

func (a *ApiService) SetPriorityChannels(
	ctx context.Context,
	req *v1.SetPriorityChannelsRequest,
) (*v1.SetPriorityChannelsResponse, error) {
	channels := req.GetPriorityChannels()

	if len(channels) == 0 {
		return nil, fmt.Errorf("priority channels list cannot be empty")
	}

	status, err := a.auth.SetPriorityChannels(ctx, channels)
	if err != nil {
		logrus.WithError(err).Error("failed to set priority channels in auth service")
		return nil, err
	}

	return &v1.SetPriorityChannelsResponse{
		Status: status,
	}, nil
}

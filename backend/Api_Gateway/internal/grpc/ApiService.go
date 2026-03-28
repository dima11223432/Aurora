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

type API interface {
	CreateGroup(ctx context.Context, name string) (int64, error)
	DeleteGroup(ctx context.Context, groupID int64) error
	UpdateGroup(ctx context.Context, groupID int64, newGroupName string) error
	GetAllGroups(ctx context.Context) ([]models.Group, error)
	CreateContact(ctx context.Context, email string, groupID int64) (int64, error)
	UpdateContact(ctx context.Context, contactID int64, newEmail string) error
	GetAllContactsByGroupID(ctx context.Context, groupID int64) ([]models.Contact, error)
	DeleteContact(ctx context.Context, contactID int64) error
	CreateNotification(ctx context.Context, title, text string, groupID int64) (int64, error)
	SendNotification(ctx context.Context, notificationID int64) (string, error)
}

type ApiService struct {
	v1.UnimplementedApiServiceServer
	api  API
	auth *services.AuthService
}

func RegisterGrpcServer(gRPC *grpc.Server, api API, auth *services.AuthService) {
	v1.RegisterApiServiceServer(gRPC, &ApiService{
		auth: auth,
		api:  api,
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

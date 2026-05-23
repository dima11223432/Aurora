// Package grpcauth implements the gRPC server for the SSO auth service.
package grpcauth

import (
	"authService/internal/domain/models"
	"authService/internal/services/auth"
	"authService/internal/storage"
	"context"
	"errors"

	ssov1 "authService/api/gen/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Auth defines the authentication operations exposed via gRPC.
type Auth interface {
	Login(ctx context.Context, user models.User, appId int) (token string, err error)
	IsAdmin(ctx context.Context, telegram_id int64) (bool, error)
	SetPriorityChannels(ctx context.Context, user_id int64, channels []string) error
	DeletePriorityChannels(ctx context.Context, user_id int64, channels []string) error
}

type serverAPI struct {
	ssov1.UnimplementedAuthServiceServer
	auth Auth
}

const (
	emptyValue = 0
)

// Register registers the AuthService gRPC server with the given gRPC server.
func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServiceServer(gRPC, &serverAPI{
		auth: auth,
	})
}

// NewServerAPI creates a new serverAPI instance.
func NewServerAPI(auth Auth) *serverAPI {
	return &serverAPI{
		auth: auth,
	}
}

// Login handles user login requests via gRPC.
func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {

	if err := validatelogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, models.User{
		Telegram_id: req.GetTelegramId(),
		Username:    req.GetUsername(),
		First_name:  req.GetFirstName(),
		Last_name:   req.GetLastName(),
	}, int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {

			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.LoginResponse{
		Token: token,
	}, nil

}

// SetPriorityChannels handles priority channel setting via gRPC.
func (s *serverAPI) SetPriorityChannels(
	ctx context.Context,
	req *ssov1.SetPriorityChannelsRequest) (
	*ssov1.SetPriorityChannelsResponse, error) {

	err := s.auth.SetPriorityChannels(ctx, req.GetUserId(), req.GetChannelsUsernames())

	if err != nil {
		if errors.Is(err, storage.ErrChannelNotFound) {
			return nil, status.Error(codes.AlreadyExists, "channel already exists")
		}
		return nil, err
	}
	return &ssov1.SetPriorityChannelsResponse{}, nil

}

// DeletePriorityChannels handles priority channel deletion via gRPC.
func (s *serverAPI) DeletePriorityChannels(ctx context.Context, req *ssov1.DeletePriorityChannelsRequest) (*ssov1.DeletePriorityChannelsResponse, error) {
	err := s.auth.DeletePriorityChannels(ctx, req.GetUserId(), req.GetChannelsUsernames())
	if err != nil {
		if errors.Is(err, storage.ErrChannelNotFound) {
			return nil, err
		}
	}
	return &ssov1.DeletePriorityChannelsResponse{}, nil
}

// IsAdmin handles admin status checks via gRPC.
func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetTelegramId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validatelogin(req *ssov1.LoginRequest) error {
	if req.GetTelegramId() == emptyValue {
		return status.Error(codes.InvalidArgument, "telegram_id is invalid")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetTelegramId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user is required")
	}
	return nil
}

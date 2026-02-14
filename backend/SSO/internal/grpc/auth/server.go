package grpcauth

import (
	"authService/internal/services/auth"
	"authService/internal/storage"
	"context"
	"errors"

	ssov1 "github.com/dima11223432/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appId int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string, isAdmin bool) (UserId int64, err error)
	IsAdmin(ctx context.Context, UserId int) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

const (
	emptyValue = 0
)

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{
		auth: auth,
	})
}
func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {

	if err := validatelogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
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

// Register register new user
func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := valdateRegister(req); err != nil {
		return nil, err
	}
	UserID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword(), req.GetIsAdmin())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.RegisterResponse{
		UserId: int64(UserID),
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}
	isAdmin, err := s.auth.IsAdmin(ctx, int(req.GetUserId()))
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
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func valdateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}
	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user is required")
	}
	return nil
}

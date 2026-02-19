package grpc

import (
	v1 "API_Service/api/gen/v1"
	services "API_Service/internal/services"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ApiService struct {
	v1.UnimplementedApiServiceServer
	auth *services.AuthService
}

func RegisterGrpcServer(gRPC *grpc.Server, auth *services.AuthService) {
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

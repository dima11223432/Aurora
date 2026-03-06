package grpc

import (
	v1 "API_Service/api/gen/v1"

	custom_errors "API_Service/internal/custom_errors"
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	return &v1.LoginResponse{
		Token: token,
	}, nil
}

func (a *ApiService) SetPriorityChannels(
	ctx context.Context,
	req *v1.SetPriorityChannelsRequest,
) (*v1.SetPriorityChannelsResponse, error) {
	channels := req.GetPriorityChannels()

	statusCode, err := a.auth.SetPriorityChannels(ctx, channels)
	if err != nil {
		if errors.Is(err, custom_errors.ErrChannelExists) {
			return nil, status.Error(codes.AlreadyExists, "channel already exists")
		}
		return nil, status.Error(status.Code(err), err.Error())
	}

	return &v1.SetPriorityChannelsResponse{
		Status: statusCode,
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

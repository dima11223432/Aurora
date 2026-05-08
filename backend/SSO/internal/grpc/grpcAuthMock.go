package grpcauth_test

import (
	"authService/internal/domain/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type GrpcAuthMock struct {
	mock.Mock
}

func (g *GrpcAuthMock) Login(ctx context.Context, user models.User, appId int) (token string, err error) {
	args := g.Called(ctx, user, appId)
	return args.String(0), args.Error(1)
}

func (g *GrpcAuthMock) DeletePriorityChannels(ctx context.Context, user_id int64, channels []string) error {
	args := g.Called(ctx, user_id, channels)
	return args.Error(0)
}

func (g *GrpcAuthMock) SetPriorityChannels(ctx context.Context, user_id int64, channels []string) error {
	args := g.Called(ctx, user_id, channels)
	return args.Error(0)
}

func (g *GrpcAuthMock) IsAdmin(ctx context.Context, telegram_id int64) (bool, error) {
	args := g.Called(ctx, telegram_id)
	return args.Bool(0), args.Error(1)
}

package grpcauth

import (
	"context"
	"recommendationService/internal/domain/models"

	"github.com/stretchr/testify/mock"
)

type UserDataProviderMock struct {
	mock.Mock
}

func (u *UserDataProviderMock) GetUserPriorityChannels(ctx context.Context, userID int64) ([]models.PriorityChannel, error) {

	args := u.Called(ctx, userID)
	return args.Get(0).([]models.PriorityChannel), args.Error(1)
}

func (g *UserDataProviderMock) GetAllParsingChannels(ctx context.Context) ([]string, error) {
	args := g.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (g *UserDataProviderMock) AddNewParsingChannel(ctx context.Context, channel string) error {
	args := g.Called(ctx, channel)
	return args.Error(0)
}

func (g *UserDataProviderMock) DeleteParsingChannel(ctx context.Context, channel string) error {
	args := g.Called(ctx, channel)
	return args.Error(0)
}

func (g *UserDataProviderMock) GetRecommendatedPosts(ctx context.Context, userID int64, cursor *models.Cursor) ([]models.Post, *models.Cursor, error) {
	args := g.Called(ctx, userID, cursor)
	return args.Get(0).([]models.Post), args.Get(1).(*models.Cursor), args.Error(2)
}

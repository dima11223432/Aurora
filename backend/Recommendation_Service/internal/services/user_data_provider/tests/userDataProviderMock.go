package auth_test

import (
	"context"
	"recommendationService/internal/domain/models"

	"github.com/stretchr/testify/mock"
)

type userDataProviderMock struct {
	mock.Mock
}

func (u *userDataProviderMock) GetPriorityChannelsByUserID(ctx context.Context, userID int64) ([]models.PriorityChannel, error) {
	args := u.Called(ctx, userID)
	return args.Get(0).([]models.PriorityChannel), args.Error(1)
}

func (u *userDataProviderMock) GetPostsByChannels(ctx context.Context, channels []string, userID int64, cursor *models.Cursor, limit int64) ([]models.Post, *models.Cursor, error) {
	args := u.Called(ctx, channels, userID, cursor, limit)
	return args.Get(0).([]models.Post), args.Get(1).(*models.Cursor), args.Error(2)
}

func (u *userDataProviderMock) GetRecommendatedPosts(ctx context.Context, userID int64, cursor *models.Cursor) ([]models.Post, *models.Cursor, error) {
	args := u.Called(ctx, userID, cursor)
	return args.Get(0).([]models.Post), args.Get(1).(*models.Cursor), args.Error(2)
}

func (u *userDataProviderMock) GetAllParsingChannels(ctx context.Context) ([]string, error) {
	args := u.Called(ctx)

	var res []string
	if args.Get(0) != nil {
		res = args.Get(0).([]string)
	}

	return res, args.Error(1)
}

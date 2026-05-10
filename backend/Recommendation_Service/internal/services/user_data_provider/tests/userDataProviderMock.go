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

func (u *userDataProviderMock) GetAllDefaultParsingChannels(ctx context.Context) ([]string, error) {
	args := u.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (u *userDataProviderMock) AddNewDefaultParsingChannel(ctx context.Context, channel string, category string) error {
	args := u.Called(ctx, channel, category)
	return args.Error(0)
}

func (u *userDataProviderMock) DeleteDefaultParsingChannel(ctx context.Context, channel string) error {
	args := u.Called(ctx, channel)
	return args.Error(0)
}

func (u *userDataProviderMock) DeleteUserCustomParsingChannel(ctx context.Context, userID int64, channel string) error {
	args := u.Called(ctx, userID, channel)
	return args.Error(0)
}

func (u *userDataProviderMock) DeleteParsingChannel(ctx context.Context, channel string) error {
	args := u.Called(ctx, channel)
	return args.Error(0)
}

func (u *userDataProviderMock) GetAllUserCustomParsingChannels(ctx context.Context, userID int64) ([]string, error) {
	args := u.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (u *userDataProviderMock) GetAllCategories(ctx context.Context) ([]string, error) {
	args := u.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (u *userDataProviderMock) AddNewUserCustomParsingChannel(ctx context.Context, userID int64, channel string) error {
	args := u.Called(ctx, userID, channel)
	return args.Error(0)
}

func (u *userDataProviderMock) GetDefaultParsingChannelsByCategory(ctx context.Context, category string) ([]string, error) {
	args := u.Called(ctx, category)
	return args.Get(0).([]string), args.Error(1)
}

func (u *userDataProviderMock) GetRecommendatedPosts(ctx context.Context, userID int64, cursor *models.Cursor) ([]models.Post, *models.Cursor, error) {
	args := u.Called(ctx, userID, cursor)
	var posts []models.Post
	var nextCursor *models.Cursor
	if args.Get(0) != nil {
		posts = args.Get(0).([]models.Post)
	}
	if args.Get(1) != nil {
		nextCursor = args.Get(1).(*models.Cursor)
	}
	return posts, nextCursor, args.Error(2)
}

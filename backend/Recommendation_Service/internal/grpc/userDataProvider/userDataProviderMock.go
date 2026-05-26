package grpcauth

import (
	"context"
	"recommendationService/internal/domain/models"

	"github.com/stretchr/testify/mock"
)

type UserDataProviderMock struct {
	mock.Mock
}

func (m *UserDataProviderMock) GetUserPriorityChannels(ctx context.Context, userID int64) ([]models.PriorityChannel, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PriorityChannel), args.Error(1)
}

func (m *UserDataProviderMock) GetRecommendatedPosts(ctx context.Context, userID int64, cursor *models.Cursor) ([]models.Post, *models.Cursor, error) {
	args := m.Called(ctx, userID, cursor)
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

func (m *UserDataProviderMock) GetAllDefaultParsingChannels(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *UserDataProviderMock) AddNewDefaultParsingChannel(ctx context.Context, channel string, category string) error {
	args := m.Called(ctx, channel, category)
	return args.Error(0)
}

func (m *UserDataProviderMock) AddNewUserCustomParsingChannel(ctx context.Context, userID int64, channel string) error {
	args := m.Called(ctx, userID, channel)
	return args.Error(0)
}

func (m *UserDataProviderMock) DeleteDefaultParsingChannel(ctx context.Context, channel string) error {
	args := m.Called(ctx, channel)
	return args.Error(0)
}

func (m *UserDataProviderMock) DeleteUserCustomParsingChannel(ctx context.Context, userID int64, channel string) error {
	args := m.Called(ctx, userID, channel)
	return args.Error(0)
}

func (m *UserDataProviderMock) GetAllUserCustomParsingChannels(ctx context.Context, userID int64) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *UserDataProviderMock) GetDefaultParsingChannelsWithCategories(ctx context.Context) (map[string][]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string][]string), args.Error(1)
}

func (m *UserDataProviderMock) GetAllCategories(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

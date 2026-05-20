package auth

import (
	"context"
	"log/slog"
	"recommendationService/internal/domain/models"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type userDataProviderSuite struct {
	suite.Suite
	mock    *userDataProviderMock
	service *UserDataProvider
}

func (u *userDataProviderSuite) SetupTest() {
	u.mock = new(userDataProviderMock)
	u.service = New(slog.Default(), u.mock, u.mock, u.mock, 5*time.Minute)
}

func (u *userDataProviderSuite) TestGetUserPriorityChannels() {
	ctx := context.Background()
	u.mock.On("GetAllDefaultParsingChannels", ctx).Return([]string{"channel1", "channel2"}, nil)
	channels, err := u.service.GetAllDefaultParsingChannels(ctx)
	u.NoError(err)
	u.Equal([]string{"channel1", "channel2"}, channels)
}

func (u *userDataProviderSuite) TestGetUserPrioityChannels() {
	ctx := context.Background()
	u.mock.On("GetPriorityChannelsByUserID", ctx, int64(1)).Return([]models.PriorityChannel{models.PriorityChannel{Channel: "channel1"}}, nil)
	channels, err := u.service.GetUserPriorityChannels(ctx, 1)
	u.NoError(err)
	u.NotEmpty(channels)
}

func (u *userDataProviderSuite) TestGetRecommendatedPosts() {
	ctx := context.Background()
	userID := int64(1)
	u.mock.On("GetPriorityChannelsByUserID", ctx, userID).
		Return([]models.PriorityChannel{{Channel: "Kafka_Channel1"}}, nil)
	u.mock.On("GetPostsByChannels", ctx, []string{"Kafka_Channel1"}, userID, mock.Anything, int64(5)).
		Return([]models.Post{}, (*models.Cursor)(nil), nil)
	posts, cursor, err := u.service.GetRecommendatedPosts(ctx, userID, nil)
	u.NoError(err)
	u.Empty(posts)
	u.Empty(cursor)
}

func (u *userDataProviderSuite) TestDeleteParsingChannel() {
	ctx := context.Background()
	channel := "test_channel"
	u.mock.On("DeleteDefaultParsingChannel", ctx, channel).Return(nil)
	u.mock.On("DeleteParsingChannel", ctx, channel).Return(nil)
	err := u.service.DeleteDefaultParsingChannel(ctx, channel)
	u.NoError(err)
}

func TestUserDataProviderTestSuite(t *testing.T) {
	suite.Run(t, new(userDataProviderSuite))
}

package auth_test

import (
	"context"
	"log/slog"
	"recommendationService/internal/domain/models"
	auth "recommendationService/internal/services/user_data_provider"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type userDataProviderSuite struct {
	suite.Suite
	mock    *userDataProviderMock
	service *auth.UserDataProvider
}

func (u *userDataProviderSuite) SetupTest() {
	u.mock = new(userDataProviderMock)
	u.service = auth.New(slog.Default(), u.mock, u.mock, u.mock, 5*time.Minute)
}

func (u *userDataProviderSuite) TestGetUserPriorityChannels() {
	ctx := context.Background()
	u.mock.On("GetAllParsingChannels", ctx).Return([]string{"channel1", "channel2"}, nil)
	channels, err := u.service.GetAllParsingChannels(ctx)
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
	u.mock.On("GetPostsByChannels", ctx, []string{"Kafka_Channel1"}, userID, mock.Anything, int64(4)).
		Return([]models.Post{}, (*models.Cursor)(nil), nil)
	posts, cursor, err := u.service.GetRecommendatedPosts(ctx, userID, nil)
	u.NoError(err)
	u.Empty(posts)
	u.Empty(cursor)
}

func TestUserDataProviderTestSuite(t *testing.T) {
	suite.Run(t, new(userDataProviderSuite))
}

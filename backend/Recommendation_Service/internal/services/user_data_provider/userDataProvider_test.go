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

func (u *userDataProviderSuite) TestAddNewUserCustomParsingChannel() {
	ctx := context.Background()
	userID := int64(1)
	channel := "custom_channel"
	u.mock.On("AddNewUserCustomParsingChannel", ctx, userID, channel).Return(nil)
	err := u.service.AddNewUserCustomParsingChannel(ctx, userID, channel)
	u.NoError(err)
}

func (u *userDataProviderSuite) TestAddNewDefaultParsingChannel() {
	ctx := context.Background()
	channel := "new_channel"
	category := "tech"
	u.mock.On("AddNewDefaultParsingChannel", ctx, channel, category).Return(nil)
	u.mock.On("AddNewParsingChannel", ctx, channel).Return(nil)
	u.mock.On("SetChannelCategory", ctx, channel, category).Return(nil)
	err := u.service.AddNewDefaultParsingChannel(ctx, channel, category)
	u.NoError(err)
}

func (u *userDataProviderSuite) TestAddNewDefaultParsingChannelChannelExistsError() {
	ctx := context.Background()
	channel := "existing_channel"
	category := "tech"
	u.mock.On("AddNewDefaultParsingChannel", ctx, channel, category).Return(nil)
	u.mock.On("AddNewParsingChannel", ctx, channel).Return(nil)
	u.mock.On("SetChannelCategory", ctx, channel, category).Return(nil)
	err := u.service.AddNewDefaultParsingChannel(ctx, channel, category)
	u.NoError(err)
}

func (u *userDataProviderSuite) TestDeleteUserCustomParsingChannel() {
	ctx := context.Background()
	userID := int64(1)
	channel := "user_channel"
	u.mock.On("DeleteUserCustomParsingChannel", ctx, userID, channel).Return(nil)
	u.mock.On("DeleteParsingChannel", ctx, channel).Return(nil)
	err := u.service.DeleteUserCustomParsingChannel(ctx, userID, channel)
	u.NoError(err)
}

func (u *userDataProviderSuite) TestGetAllUserCustomParsingChannels() {
	ctx := context.Background()
	userID := int64(1)
	u.mock.On("GetAllUserCustomParsingChannels", ctx, userID).Return([]string{"ch1", "ch2"}, nil)
	channels, err := u.service.GetAllUserCustomParsingChannels(ctx, userID)
	u.NoError(err)
	u.Equal([]string{"ch1", "ch2"}, channels)
}

func (u *userDataProviderSuite) TestGetDefaultParsingChannelsWithCategories() {
	ctx := context.Background()
	u.mock.On("GetAllCategories", ctx).Return([]string{"tech", "news"}, nil)
	u.mock.On("GetDefaultParsingChannelsByCategory", ctx, "tech").Return([]string{"ch1", "ch2"}, nil)
	u.mock.On("GetDefaultParsingChannelsByCategory", ctx, "news").Return([]string{"ch3"}, nil)
	result, err := u.service.GetDefaultParsingChannelsWithCategories(ctx)
	u.NoError(err)
	u.Equal(map[string][]string{"tech": {"ch1", "ch2"}, "news": {"ch3"}}, result)
}

func (u *userDataProviderSuite) TestGetDefaultParsingChannelsWithCategoriesEmpty() {
	ctx := context.Background()
	u.mock.On("GetAllCategories", ctx).Return([]string{}, nil)
	result, err := u.service.GetDefaultParsingChannelsWithCategories(ctx)
	u.NoError(err)
	u.Equal(map[string][]string{}, result)
}

func (u *userDataProviderSuite) TestGetDefaultParsingChannelsByCategory() {
	ctx := context.Background()
	category := "tech"
	u.mock.On("GetDefaultParsingChannelsByCategory", ctx, category).Return([]string{"ch1", "ch2"}, nil)
	channels, err := u.service.GetDefaultParsingChannelsByCategory(ctx, category)
	u.NoError(err)
	u.Equal([]string{"ch1", "ch2"}, channels)
}

func (u *userDataProviderSuite) TestGetAllCategories() {
	ctx := context.Background()
	u.mock.On("GetAllCategories", ctx).Return([]string{"tech", "news", "sports"}, nil)
	categories, err := u.service.GetAllCategories(ctx)
	u.NoError(err)
	u.Equal([]string{"tech", "news", "sports"}, categories)
}

func TestUserDataProviderTestSuite(t *testing.T) {
	suite.Run(t, new(userDataProviderSuite))
}

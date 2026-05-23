package grpcauth

import (
	"context"
	"errors"
	recv1 "recommendationService/api/gen/v1"
	"recommendationService/internal/domain/models"
	"recommendationService/internal/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserDataProviderSuite struct {
	suite.Suite
	mock *UserDataProviderMock
	s    recv1.RecommendationServiceServer
}

func (u *UserDataProviderSuite) SetupTest() {
	u.mock = new(UserDataProviderMock)
	u.s = NewServerAPI(u.mock, u.mock, u.mock)
}

func (u *UserDataProviderSuite) TestGetUserPriorityChannels() {
	ctx := context.Background()
	userID := int64(1)
	u.mock.On("GetUserPriorityChannels", ctx, userID).
		Return([]models.PriorityChannel{models.PriorityChannel{Channel: "channel1"}}, nil)
	channels, err := u.s.GetUserPriorityChannels(ctx, &recv1.GetUserPriorityChannelsRequest{UserId: 1})
	u.NoError(err)
	u.NotNil(channels)
}

func TestUserDataProviderSuite(t *testing.T) {
	suite.Run(t, new(UserDataProviderSuite))
}

func (u *UserDataProviderSuite) TestGetAllParsingChannels() {
	ctx := context.Background()
	expected := []string{"channel1", "channel2"}

	u.mock.On("GetAllDefaultParsingChannels", ctx).Return(expected, nil)

	resp, err := u.s.GetAllDefaultParsingChannels(ctx, &recv1.GetAllDefaultParsingChannelsRequest{})

	u.NoError(err)
	u.NotNil(resp)
	u.Equal(expected, resp.Channels)
	u.mock.AssertExpectations(u.T())
}

func (u *UserDataProviderSuite) TestGetRecommendatedPosts_Success() {
	ctx := context.Background()
	userID := int64(1)

	mockNextCursor := &models.Cursor{Score: 1.5, ID: "test_id"}

	u.mock.On("GetRecommendatedPosts", ctx, userID, (*models.Cursor)(nil)).
		Return([]models.Post{
			{
				PostText: "Test text",
				Date:     time.Now(),
			},
		}, mockNextCursor, nil)

	resp, err := u.s.GetRecommendatedPosts(ctx, &recv1.GetRecommendatedPostsRequest{UserId: userID})

	u.NoError(err)
	u.NotNil(resp)
	u.Len(resp.Posts, 1)
	u.Equal("Test text", resp.Posts[0].PostText)
	u.Equal(mockNextCursor.ID, resp.NextCursor.Id)
}

func (u *UserDataProviderSuite) TestDeleteParsingChannel() {
	ctx := context.Background()
	channelUsername := "channel1"

	u.mock.On("DeleteDefaultParsingChannel", ctx, channelUsername).Return(nil).Once()

	resp, err := u.s.DeleteDefaultParsingChannel(ctx, &recv1.DeleteDefaultParsingChannelRequest{
		ChannelUsername: channelUsername,
	})

	u.NoError(err)
	u.NotNil(resp)
	u.mock.AssertExpectations(u.T())

	resp, err = u.s.DeleteDefaultParsingChannel(ctx, &recv1.DeleteDefaultParsingChannelRequest{
		ChannelUsername: "",
	})

	u.Nil(resp)
	u.Error(err)
	u.Equal(codes.InvalidArgument, status.Code(err))
	u.mock.AssertNotCalled(u.T(), "DeleteDefaultParsingChannel", ctx, "")
}

func (u *UserDataProviderSuite) TestGetAllUserCustomParsingChannels() {
	ctx := context.Background()
	userID := int64(1)
	expected := []string{"user_ch1", "user_ch2"}

	u.mock.On("GetAllUserCustomParsingChannels", ctx, userID).Return(expected, nil)

	resp, err := u.s.GetAllUserCustomParsingChannels(ctx, &recv1.GetAllUserCustomParsingChannelsRequest{UserId: userID})

	u.NoError(err)
	u.NotNil(resp)
	u.Equal(expected, resp.Channels)
}

func (u *UserDataProviderSuite) TestGetAllUserCustomParsingChannelsError() {
	ctx := context.Background()
	userID := int64(1)

	u.mock.On("GetAllUserCustomParsingChannels", ctx, userID).Return(nil, errors.New("db error"))

	_, err := u.s.GetAllUserCustomParsingChannels(ctx, &recv1.GetAllUserCustomParsingChannelsRequest{UserId: userID})

	u.Error(err)
	u.Equal(codes.Internal, status.Code(err))
}

func (u *UserDataProviderSuite) TestAddNewUserCustomParsingChannel() {
	ctx := context.Background()
	userID := int64(1)
	channel := "new_channel"

	u.mock.On("AddNewUserCustomParsingChannel", ctx, userID, channel).Return(nil)

	resp, err := u.s.AddNewUserCustomParsingChannel(ctx, &recv1.AddNewUserCustomParsingChannelRequest{
		UserId:          userID,
		ChannelUsername: channel,
	})

	u.NoError(err)
	u.NotNil(resp)
}

func (u *UserDataProviderSuite) TestAddNewUserCustomParsingChannelAlreadyExists() {
	ctx := context.Background()
	userID := int64(1)
	channel := "existing_channel"

	u.mock.On("AddNewUserCustomParsingChannel", ctx, userID, channel).Return(storage.ErrChannelExists)

	_, err := u.s.AddNewUserCustomParsingChannel(ctx, &recv1.AddNewUserCustomParsingChannelRequest{
		UserId:          userID,
		ChannelUsername: channel,
	})

	u.Error(err)
	u.Equal(codes.AlreadyExists, status.Code(err))
}

func (u *UserDataProviderSuite) TestAddNewUserCustomParsingChannelError() {
	ctx := context.Background()
	userID := int64(1)
	channel := "new_channel"

	u.mock.On("AddNewUserCustomParsingChannel", ctx, userID, channel).Return(errors.New("db error"))

	_, err := u.s.AddNewUserCustomParsingChannel(ctx, &recv1.AddNewUserCustomParsingChannelRequest{
		UserId:          userID,
		ChannelUsername: channel,
	})

	u.Error(err)
	u.Equal(codes.Internal, status.Code(err))
}

func (u *UserDataProviderSuite) TestAddNewDefaultParsingChannel() {
	ctx := context.Background()
	channel := "new_channel"
	category := "tech"

	u.mock.On("AddNewDefaultParsingChannel", ctx, channel, category).Return(nil)

	resp, err := u.s.AddNewDefaultParsingChannel(ctx, &recv1.AddNewDefaultParsingChannelRequest{
		ChannelUsername: channel,
		Category:        category,
	})

	u.NoError(err)
	u.NotNil(resp)
}

func (u *UserDataProviderSuite) TestAddNewDefaultParsingChannelEmpty() {
	ctx := context.Background()

	resp, err := u.s.AddNewDefaultParsingChannel(ctx, &recv1.AddNewDefaultParsingChannelRequest{
		ChannelUsername: "",
		Category:        "tech",
	})

	u.Nil(resp)
	u.Error(err)
	u.Equal(codes.InvalidArgument, status.Code(err))
}

func (u *UserDataProviderSuite) TestAddNewDefaultParsingChannelAlreadyExists() {
	ctx := context.Background()
	channel := "existing_channel"
	category := "tech"

	u.mock.On("AddNewDefaultParsingChannel", ctx, channel, category).Return(storage.ErrChannelExists)

	_, err := u.s.AddNewDefaultParsingChannel(ctx, &recv1.AddNewDefaultParsingChannelRequest{
		ChannelUsername: channel,
		Category:        category,
	})

	u.Error(err)
	u.Equal(codes.AlreadyExists, status.Code(err))
}

func (u *UserDataProviderSuite) TestAddNewDefaultParsingChannelError() {
	ctx := context.Background()
	channel := "new_channel"
	category := "tech"

	u.mock.On("AddNewDefaultParsingChannel", ctx, channel, category).Return(errors.New("db error"))

	_, err := u.s.AddNewDefaultParsingChannel(ctx, &recv1.AddNewDefaultParsingChannelRequest{
		ChannelUsername: channel,
		Category:        category,
	})

	u.Error(err)
	u.Equal(codes.Internal, status.Code(err))
}

func (u *UserDataProviderSuite) TestDeleteUserCustomParsingChannel() {
	ctx := context.Background()
	userID := int64(1)
	channel := "user_channel"

	u.mock.On("DeleteUserCustomParsingChannel", ctx, userID, channel).Return(nil)

	resp, err := u.s.DeleteUserCustomParsingChannel(ctx, &recv1.DeleteUserCustomParsingChannelRequest{
		UserId:          userID,
		ChannelUsername: channel,
	})

	u.NoError(err)
	u.NotNil(resp)
}

func (u *UserDataProviderSuite) TestDeleteUserCustomParsingChannelError() {
	ctx := context.Background()
	userID := int64(1)
	channel := "user_channel"

	u.mock.On("DeleteUserCustomParsingChannel", ctx, userID, channel).Return(errors.New("db error"))

	_, err := u.s.DeleteUserCustomParsingChannel(ctx, &recv1.DeleteUserCustomParsingChannelRequest{
		UserId:          userID,
		ChannelUsername: channel,
	})

	u.Error(err)
	u.Equal(codes.Internal, status.Code(err))
}

func (u *UserDataProviderSuite) TestGetAllDefaultParsingChannelsWithCategories() {
	ctx := context.Background()
	expected := map[string][]string{
		"tech":  {"ch1", "ch2"},
		"news":  {"ch3"},
		"sports": {"ch4"},
	}

	u.mock.On("GetDefaultParsingChannelsWithCategories", ctx).Return(expected, nil)

	resp, err := u.s.GetAllDefaultParsingChannelsWithCategories(ctx, &recv1.GetAllDefaultParsingChannelsWithCategoriesRequest{})

	u.NoError(err)
	u.NotNil(resp)
	u.Len(resp.Channels, 3)
	u.Equal([]string{"ch1", "ch2"}, resp.Channels["tech"].Usernames)
}

func (u *UserDataProviderSuite) TestGetAllDefaultParsingChannelsWithCategoriesError() {
	ctx := context.Background()

	u.mock.On("GetDefaultParsingChannelsWithCategories", ctx).Return(nil, errors.New("db error"))

	_, err := u.s.GetAllDefaultParsingChannelsWithCategories(ctx, &recv1.GetAllDefaultParsingChannelsWithCategoriesRequest{})

	u.Error(err)
	u.Equal(codes.Internal, status.Code(err))
}

func (u *UserDataProviderSuite) TestGetUserPriorityChannelsError() {
	ctx := context.Background()
	userID := int64(1)

	u.mock.On("GetUserPriorityChannels", ctx, userID).Return(nil, errors.New("db error"))

	_, err := u.s.GetUserPriorityChannels(ctx, &recv1.GetUserPriorityChannelsRequest{UserId: userID})

	u.Error(err)
	u.Equal(codes.Internal, status.Code(err))
}

func (u *UserDataProviderSuite) TestGetAllDefaultParsingChannelsError() {
	ctx := context.Background()

	u.mock.On("GetAllDefaultParsingChannels", ctx).Return(nil, errors.New("db error"))

	_, err := u.s.GetAllDefaultParsingChannels(ctx, &recv1.GetAllDefaultParsingChannelsRequest{})

	u.Error(err)
	u.Equal(codes.Internal, status.Code(err))
}

func (u *UserDataProviderSuite) TestGetRecommendatedPostsError() {
	ctx := context.Background()
	userID := int64(1)

	mockCursor := &models.Cursor{Score: 1.5, ID: "test_id"}
	u.mock.On("GetRecommendatedPosts", ctx, userID, (*models.Cursor)(nil)).Return([]models.Post{}, mockCursor, errors.New("db error"))

	_, err := u.s.GetRecommendatedPosts(ctx, &recv1.GetRecommendatedPostsRequest{UserId: userID})

	u.Error(err)
	u.Equal(codes.Internal, status.Code(err))
}

func (u *UserDataProviderSuite) TestGetRecommendatedPostsWithCursor() {
	ctx := context.Background()
	userID := int64(1)

	reqCursor := &recv1.Cursor{Score: 1, Id: "cursor_id"}
	mockCursor := &models.Cursor{Score: 1, ID: "cursor_id"}
	u.mock.On("GetRecommendatedPosts", ctx, userID, mockCursor).Return([]models.Post{}, mockCursor, nil)

	resp, err := u.s.GetRecommendatedPosts(ctx, &recv1.GetRecommendatedPostsRequest{UserId: userID, Cursor: reqCursor})

	u.NoError(err)
	u.NotNil(resp)
}

func (u *UserDataProviderSuite) TestDeleteDefaultParsingChannelError() {
	ctx := context.Background()
	channelUsername := "channel1"

	u.mock.On("DeleteDefaultParsingChannel", ctx, channelUsername).Return(errors.New("db error"))

	_, err := u.s.DeleteDefaultParsingChannel(ctx, &recv1.DeleteDefaultParsingChannelRequest{
		ChannelUsername: channelUsername,
	})

	u.Error(err)
	u.Equal(codes.Internal, status.Code(err))
}

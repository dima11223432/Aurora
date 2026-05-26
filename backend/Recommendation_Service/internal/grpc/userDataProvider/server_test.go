package grpcauth

import (
	"context"
	recv1 "recommendationService/api/gen/v1"
	"recommendationService/internal/domain/models"
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

func (u *UserDataProviderSuite) TestAddNewParsingChannel() {
    ctx := context.Background()
    
    u.Run("Success", func() {
        channelUsername := "channel"
        req := &recv1.AddNewParsingChannelRequest{
            ChannelUsername: channelUsername,
        }
        
        u.mock.On("AddNewParsingChannel", ctx, channelUsername).
            Return(nil).Once()
        
        resp, err := u.s.AddNewParsingChannel(ctx, req)
        
        u.NoError(err)
        u.NotNil(resp)
        u.mock.AssertExpectations(u.T())
    })

	u.Run("Invalid Input", func() {
        req := &recv1.AddNewParsingChannelRequest{
            ChannelUsername: "", 
        }
        
        resp, err := u.s.AddNewParsingChannel(ctx, req)
        
        u.Error(err) 
        
        status, ok := status.FromError(err)
        u.True(ok, "Должна быть ошибка")
        u.Equal(codes.InvalidArgument, status.Code(), "Неверный код")
        
        u.Nil(resp)
        
        u.mock.AssertNotCalled(u.T(), "AddNewParsingChannel", ctx, "")
    })
    
	u.Run("Database Error", func() {
        channelUsername := "channel"
        req := &recv1.AddNewParsingChannelRequest{
            ChannelUsername: channelUsername,
        }
        
        dbError := errors.New("db error")
        u.mock.On("AddNewParsingChannel", ctx, channelUsername).
            Return(dbError).Once()
        
        resp, err := u.s.AddNewParsingChannel(ctx, req)
        
        u.Error(err)
        u.Nil(resp)

		u.mock.AssertExpectations(u.T())
	})

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

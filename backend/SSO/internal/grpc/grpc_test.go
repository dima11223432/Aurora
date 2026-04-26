package grpcauth_test

import (
	ssov1 "authService/api/gen/v1"
	"authService/internal/grpc/auth"
	"authService/internal/storage"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServerTestSuite struct {
	suite.Suite
	authMock *GrpcAuthMock
	server   ssov1.AuthServiceServer
}

func (s *AuthServerTestSuite) SetupTest() {
	s.authMock = new(GrpcAuthMock)
	s.server = grpcauth.NewServerAPI(s.authMock)
}

func (s *AuthServerTestSuite) TestSetPriorityChannels_ConflictMapping() {
	userID := int64(1)
	channels := []string{"channel1"}

	s.authMock.On("SetPriorityChannels", mock.Anything, userID, channels).
		Return(storage.ErrChannelNotFound).Once()

	req := &ssov1.SetPriorityChannelsRequest{
		UserId:            userID,
		ChannelsUsernames: channels,
	}

	resp, err := s.server.SetPriorityChannels(context.Background(), req)

	s.Nil(resp)
	s.Error(err)
	s.Equal(codes.AlreadyExists, status.Code(err))
	s.Equal("channel already exists", status.Convert(err).Message())
	s.authMock.AssertExpectations(s.T())
}

func (s *AuthServerTestSuite) TestDeletePriorityChannels_SuccessReturnsEmptyResponse() {
	userID := int64(1)
	channels := []string{"channel1", "channel2"}

	s.authMock.On("DeletePriorityChannels", mock.Anything, userID, channels).
		Return(nil).Once()

	req := &ssov1.DeletePriorityChannelsRequest{
		UserId:            userID,
		ChannelsUsernames: channels,
	}
	resp, err := s.server.DeletePriorityChannels(context.Background(), req)

	s.NoError(err)
	s.Equal(&ssov1.DeletePriorityChannelsResponse{}, resp)
	s.authMock.AssertExpectations(s.T())
}

func (s *AuthServerTestSuite) TestDeletePriorityChannels_ServiceErrorPropagates() {
	userID := int64(1)
	channels := []string{"channel1", "channel2"}

	s.authMock.On("DeletePriorityChannels", mock.Anything, userID, channels).
		Return(storage.ErrChannelNotFound).Once()

	req := &ssov1.DeletePriorityChannelsRequest{
		UserId:            userID,
		ChannelsUsernames: channels,
	}
	resp, err := s.server.DeletePriorityChannels(context.Background(), req)

	s.Nil(resp)
	s.Error(err)
	s.ErrorIs(err, storage.ErrChannelNotFound)
	s.authMock.AssertExpectations(s.T())
}

func (s *AuthServerTestSuite) TestSetPriorityChannels_EmptySlice_NoPanic() {
	userID := int64(1)
	channels := []string{}

	s.authMock.On("SetPriorityChannels", mock.Anything, userID, channels).
		Return(nil).Once()

	req := &ssov1.SetPriorityChannelsRequest{
		UserId:            userID,
		ChannelsUsernames: channels,
	}

	resp, err := s.server.SetPriorityChannels(context.Background(), req)

	s.NoError(err)
	s.Equal(&ssov1.SetPriorityChannelsResponse{}, resp)
	s.authMock.AssertExpectations(s.T())
}

func (s *AuthServerTestSuite) TestDeletePriorityChannels_EmptySlice_NoPanic() {
	userID := int64(1)
	channels := []string{}

	s.authMock.On("DeletePriorityChannels", mock.Anything, userID, channels).
		Return(nil).Once()

	req := &ssov1.DeletePriorityChannelsRequest{
		UserId:            userID,
		ChannelsUsernames: channels,
	}

	resp, err := s.server.DeletePriorityChannels(context.Background(), req)

	s.NoError(err)
	s.Equal(&ssov1.DeletePriorityChannelsResponse{}, resp)
	s.authMock.AssertExpectations(s.T())
}

func TestAuthServer(t *testing.T) {
	suite.Run(t, new(AuthServerTestSuite))
}

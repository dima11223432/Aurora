package grpcauth_test

import (
	ssov1 "authService/api/gen/v1"
	"authService/internal/grpc/auth"
	"authService/internal/storage"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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
func (s *AuthServerTestSuite) TestDeletePriorityChannels_Success() {
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
	s.NotNil(resp)
	s.authMock.AssertExpectations(s.T())
}

func (s *AuthServerTestSuite) TestDeletePriorityChannels_Failure() {
	userID := int64(1)
	channels := []string{"channel1", "channel2"}

	s.authMock.On("DeletePriorityChannels", mock.Anything, userID, channels).
		Return(storage.ErrChannelNotFound).Once()

	req := &ssov1.DeletePriorityChannelsRequest{
		UserId:            userID,
		ChannelsUsernames: channels,
	}
	_, err := s.server.DeletePriorityChannels(context.Background(), req)

	s.Error(err)
	s.authMock.AssertExpectations(s.T())
}

func TestAuthServer(t *testing.T) {
	suite.Run(t, new(AuthServerTestSuite))
}

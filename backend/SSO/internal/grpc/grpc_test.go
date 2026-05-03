package grpcauth_test

import (
	ssov1 "authService/api/gen/v1"
	"authService/internal/domain/models"
	"authService/internal/grpc/auth"
	servicesauth "authService/internal/services/auth"
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

func (s *AuthServerTestSuite) TestLogin_ValidationAppIDRequired() {
	req := &ssov1.LoginRequest{
		TelegramId: 123456789,
		AppId:      0,
	}

	resp, err := s.server.Login(context.Background(), req)

	s.Nil(resp)
	s.Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Equal("app_id is required", status.Convert(err).Message())
}

func (s *AuthServerTestSuite) TestLogin_ValidationTelegramIDInvalid() {
	req := &ssov1.LoginRequest{
		TelegramId: 0,
		AppId:      1,
	}

	resp, err := s.server.Login(context.Background(), req)

	s.Nil(resp)
	s.Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Equal("telegram_id is invalid", status.Convert(err).Message())
}

func (s *AuthServerTestSuite) TestLogin_Success_MapsFieldsAndToken() {
	req := &ssov1.LoginRequest{
		TelegramId: 123456789,
		Username:   "john_doe",
		FirstName:  "John",
		LastName:   "Doe",
		AppId:      int64(7),
	}

	s.authMock.On(
		"Login",
		mock.Anything,
		mock.MatchedBy(func(user models.User) bool {
			return user.Telegram_id == req.GetTelegramId() &&
				user.Username == req.GetUsername() &&
				user.First_name == req.GetFirstName() &&
				user.Last_name == req.GetLastName()
		}),
		int(req.GetAppId()),
	).Return("jwt-token", nil).Once()

	resp, err := s.server.Login(context.Background(), req)

	s.NoError(err)
	s.NotNil(resp)
	s.Equal("jwt-token", resp.GetToken())
	s.authMock.AssertExpectations(s.T())
}

func (s *AuthServerTestSuite) TestLogin_InvalidCredentials_ReturnsInvalidArgument() {
	req := &ssov1.LoginRequest{
		TelegramId: 1001,
		Username:   "john_doe",
		FirstName:  "John",
		LastName:   "Doe",
		AppId:      1,
	}

	s.authMock.On(
		"Login",
		mock.Anything,
		mock.Anything,
		int(req.GetAppId()),
	).Return("", servicesauth.ErrInvalidCredentials).Once()

	resp, err := s.server.Login(context.Background(), req)

	s.Nil(resp)
	s.Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Equal("invalid credentials", status.Convert(err).Message())
	s.authMock.AssertExpectations(s.T())
}

func TestAuthServer(t *testing.T) {
	suite.Run(t, new(AuthServerTestSuite))
}

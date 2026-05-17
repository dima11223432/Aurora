package auth_test

import (
	"authService/internal/domain/models"
	"authService/internal/services/auth"
	"authService/internal/storage"
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

type MockUserSaver struct {
	mock.Mock
}

func (m *MockUserSaver) SaveUser(ctx context.Context, user models.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

type MockUserProvider struct {
	mock.Mock
}

func (m *MockUserProvider) User(ctx context.Context, telegramID int64) (models.User, error) {
	args := m.Called(ctx, telegramID)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserProvider) IsAdmin(ctx context.Context, telegramID int64) (bool, error) {
	args := m.Called(ctx, telegramID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserProvider) SetPriorityChannels(ctx context.Context, userID int64, channels []string) error {
	args := m.Called(ctx, userID, channels)
	return args.Error(0)
}

func (m *MockUserProvider) DeletePriorityChannels(ctx context.Context, userID int64, channels []string) error {
	args := m.Called(ctx, userID, channels)
	return args.Error(0)
}

type MockAppProvider struct {
	mock.Mock
}

func (m *MockAppProvider) App(ctx context.Context, appID int64) (models.App, error) {
	args := m.Called(ctx, appID)
	return args.Get(0).(models.App), args.Error(1)
}

type AuthTestSuite struct {
	suite.Suite

	storage *postgres.Storage
	service *auth.Auth
}

func (s *AuthTestSuite) SetupTest() {
	s.mockUserSaver = new(MockUserSaver)
	s.mockUserProvider = new(MockUserProvider)
	s.mockAppProvider = new(MockAppProvider)
	s.service = auth.New(slog.Default(), s.mockUserSaver, s.mockUserProvider, s.mockAppProvider, 5*time.Minute)
}

func (s *AuthTestSuite) TestLogin_Success() {
	app := models.App{ID: 1, Name: "test_app", Secret: "test_secret"}
	user := models.User{
		Telegram_id: 123456789,
		First_name:  "Dima",
		Last_name:   "Dmitriev",
		Username:    "dimadmitriev",
		Is_admin:    false,
	}
	savedUser := models.User{ID: 1, Telegram_id: 123456789, First_name: "Dima", Last_name: "Dmitriev", Username: "dimadmitriev", Is_admin: false}

	s.mockUserProvider.On("User", mock.Anything, int64(123456789)).Return(savedUser, nil)
	s.mockAppProvider.On("App", mock.Anything, int64(1)).Return(app, nil)

	token, err := s.service.Login(context.Background(), user, 1)
	s.NoError(err)
	s.NotEmpty(token)

	s.mockUserProvider.AssertExpectations(s.T())
	s.mockAppProvider.AssertExpectations(s.T())
}

func (s *AuthTestSuite) TestLogin_InvalidTelegramID() {
	user := models.User{
		Telegram_id: 0,
		First_name:  "Dima",
	}

	regularUser := models.User{
		Telegram_id: 222222222,
		First_name:  "Regular",
		Last_name:   "User",
		Username:    "regular_user",
		Is_admin:    false,
	}

	_, err := s.service.Login(ctx, adminUser, 1)
	s.NoError(err)

	_, err = s.service.Login(ctx, regularUser, 1)
	s.NoError(err)

	isAdmin, err := s.service.IsAdmin(ctx, adminUser.Telegram_id)
	s.NoError(err)
	s.True(isAdmin)

	isAdmin, err = s.service.IsAdmin(ctx, regularUser.Telegram_id)
	s.NoError(err)
	s.False(isAdmin)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) TestLogin_AppNotFound() {
	user := models.User{
		Telegram_id: 123456789,
		First_name:  "Dima",
		Last_name:   "Dmitriev",
		Username:    "dimadmitriev",
	}
	savedUser := models.User{ID: 1, Telegram_id: 123456789, First_name: "Dima", Last_name: "Dmitriev", Username: "dimadmitriev"}

	s.mockUserProvider.On("User", mock.Anything, int64(123456789)).Return(savedUser, nil)
	s.mockAppProvider.On("App", mock.Anything, int64(999)).Return(models.App{}, fmt.Errorf("storage.postgres.App: %w", storage.ErrAppNotFound))

	_, err := s.service.Login(context.Background(), user, 999)
	s.Error(err)
	s.ErrorContains(err, "app not found")
}

func (s *AuthTestSuite) TestRegisterNewUser_Success() {
	user := models.User{
		Telegram_id: 999999999,
		First_name:  "Test",
		Last_name:   "User",
		Username:    "testuser",
		Is_admin:    false,
	}

	s.mockUserProvider.On("User", mock.Anything, int64(999999999)).Return(models.User{}, fmt.Errorf("storage.postgres.User: %w", storage.ErrUserNotFound))
	s.mockUserSaver.On("SaveUser", mock.Anything, user).Return(int64(1), nil)

	id, err := s.service.RegisterNewUser(context.Background(), user)
	s.NoError(err)
	s.Equal(int64(1), id)

	s.mockUserProvider.AssertExpectations(s.T())
	s.mockUserSaver.AssertExpectations(s.T())
}

func (s *AuthTestSuite) TestRegisterNewUser_AlreadyExists() {
	user := models.User{
		Telegram_id: 999999999,
		First_name:  "Test",
		Last_name:   "User",
		Username:    "testuser",
		Is_admin:    false,
	}

	existingUser := models.User{ID: 1, Telegram_id: 999999999}
	s.mockUserProvider.On("User", mock.Anything, int64(999999999)).Return(existingUser, nil)

	_, err := s.service.RegisterNewUser(context.Background(), user)
	s.Error(err)
	s.ErrorIs(err, auth.ErrUserExists)
}

func (s *AuthTestSuite) TestDeletePriorityChannels() {
	s.mockUserProvider.On("DeletePriorityChannels", mock.Anything, int64(1), []string{"channel1", "channel2"}).Return(nil)

	err := s.service.DeletePriorityChannels(context.Background(), 1, []string{"channel1", "channel2"})
	s.NoError(err)

	s.mockUserProvider.AssertExpectations(s.T())
}

func (s *AuthTestSuite) TestSetPriorityChannels() {
	s.mockUserProvider.On("SetPriorityChannels", mock.Anything, int64(1), []string{"channel1", "channel2"}).Return(nil)

	err := s.service.SetPriorityChannels(context.Background(), 1, []string{"channel1", "channel2"})
	s.NoError(err)

	s.mockUserProvider.AssertExpectations(s.T())
}

func (s *AuthTestSuite) TestSetPriorityChannels_Error() {
	s.mockUserProvider.On("SetPriorityChannels", mock.Anything, int64(999999), []string{"channel1"}).Return(assert.AnError)

	err := s.service.SetPriorityChannels(context.Background(), 999999, []string{"channel1"})
	s.Error(err)

	s.mockUserProvider.AssertExpectations(s.T())
}

func (s *AuthTestSuite) TestIsAdmin_Success() {
	s.mockUserProvider.On("IsAdmin", mock.Anything, int64(111111111)).Return(true, nil)

	isAdmin, err := s.service.IsAdmin(context.Background(), 111111111)
	s.NoError(err)
	s.True(isAdmin)

	s.mockUserProvider.AssertExpectations(s.T())
}

func (s *AuthTestSuite) TestIsAdmin_UserNotFound() {
	s.mockUserProvider.On("IsAdmin", mock.Anything, int64(999999999)).Return(false, auth.ErrInvalidCredentials)

	_, err := s.service.IsAdmin(context.Background(), 999999999)
	s.Error(err)
	s.ErrorIs(err, auth.ErrInvalidCredentials)

	s.mockUserProvider.AssertExpectations(s.T())
}

func (s *AuthTestSuite) TestIsAdmin_InternalError() {
	s.mockUserProvider.On("IsAdmin", mock.Anything, int64(111111111)).Return(false, assert.AnError)

	_, err := s.service.IsAdmin(context.Background(), 111111111)
	s.Error(err)

	s.mockUserProvider.AssertExpectations(s.T())
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
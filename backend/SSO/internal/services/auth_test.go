package auth_test

import (
	"authService/internal/domain/models"
	"authService/internal/services/auth"
	"authService/internal/storage/postgres"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	storage *postgres.Storage
	service *auth.Auth
}

func (s *AuthTestSuite) SetupTest() {
	storage, err := postgres.New("postgres://postgres:pass@localhost:5432/test_auth?sslmode=disable")
	s.Require().NoError(err)

	s.storage = storage
	s.service = auth.New(slog.Default(), storage, storage, storage, 5*time.Minute)

	s.storage.DB.Exec("INSERT INTO apps (id, name, secret) VALUES (1, 'test', 'secret')")
}

func (s *AuthTestSuite) TearDownTest() {
	s.storage.DB.Exec("TRUNCATE TABLE users, channels, apps RESTART IDENTITY CASCADE")
	s.storage.DB.Close()
}

func (s *AuthTestSuite) TestLoginUser() {
	testCases := []struct {
		name    string
		user    models.User
		withErr bool
	}{
		{
			name: "Valid user",
			user: models.User{
				Telegram_id: 123456789,
				First_name:  "Dima",
				Last_name:   "Dmitriev",
				Username:    "dimadmitriev",
				Is_admin:    false,
			},
			withErr: false,
		},
		{
			name: "Invalid Telegram ID",
			user: models.User{
				Telegram_id: 0,
				First_name:  "Dima",
			},
			withErr: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			token, err := s.service.Login(context.Background(), tc.user, 1)
			if tc.withErr {
				s.Error(err)
				s.Empty(token)
			} else {
				s.NoError(err)
				s.NotEmpty(token)
			}
		})
	}
}

func (s *AuthTestSuite) TestSetPriorityChannels() {
	user := models.User{
		Telegram_id: 123456789,
		First_name:  "Dima",
		Last_name:   "Dmitriev",
		Username:    "dimadmitriev",
		Is_admin:    false,
	}

	token, err := s.service.Login(context.Background(), user, 1)
	s.NoError(err)
	s.NotEmpty(token)

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	claims := jwt.MapClaims{}
	_, err = parser.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		return []byte("secret"), nil
	})

	id := int64(claims["id"].(float64))
	s.NoError(err)
	s.NotZero(claims["id"])

	err = s.service.SetPriorityChannels(context.Background(), id, []string{"channel1", "channel2"})
	s.NoError(err)
}

func (s *AuthTestSuite) TestDeletePriorityChannels() {
	user := models.User{

		Telegram_id: 123456789,
		First_name:  "Dima",
		Last_name:   "Dmitriev",
		Username:    "dimadmitriev",
		Is_admin:    false,
	}
	token, err := s.service.Login(context.Background(), user, 1)
	s.NoError(err)
	s.NotEmpty(token)
	err = s.service.SetPriorityChannels(context.Background(), 1, []string{"channel1", "channel2"})
	s.NoError(err)
	err = s.service.DeletePriorityChannels(context.Background(), 1, []string{"channel1", "channel2"})
	s.NoError(err)
}

func (s *AuthTestSuite) TestIsAdmin() {
	ctx := context.Background()

	adminUser := models.User{
		Telegram_id: 111111111,
		First_name:  "Admin",
		Last_name:   "User",
		Username:    "admin_user",
		Is_admin:    true,
	}

	regularUser := models.User{
		Telegram_id: 222222222,
		First_name:  "Regular",
		Last_name:   "User",
		Username:    "regular_user",
		Is_admin:    false,
	}

	_, err := s.service.Login(ctx, adminUser, 1)
	assert.NoError(s.T(), err)

	_, err = s.service.Login(ctx, regularUser, 1)
	assert.NoError(s.T(), err)

	isAdmin, err := s.service.IsAdmin(ctx, adminUser.Telegram_id)
	assert.NoError(s.T(), err)
	assert.True(s.T(), isAdmin)

	isAdmin, err = s.service.IsAdmin(ctx, regularUser.Telegram_id)
	assert.NoError(s.T(), err)
	assert.False(s.T(), isAdmin)
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
		Is_admin:    false,
	}

	_, err := s.service.Login(context.Background(), user, 999)
	s.Error(err)
}

func (s *AuthTestSuite) TestRegisterNewUser_AlreadyExists() {
	ctx := context.Background()
	user := models.User{
		Telegram_id: 999999999,
		First_name:  "Test",
		Last_name:   "User",
		Username:    "testuser",
		Is_admin:    false,
	}

	id, err := s.service.RegisterNewUser(ctx, user)
	s.NoError(err)
	s.NotZero(id)

	_, err = s.service.RegisterNewUser(ctx, user)
	s.Error(err)
	s.ErrorIs(err, auth.ErrUserExists)
}

func (s *AuthTestSuite) TestIsAdmin_UserNotFound() {
	_, err := s.service.IsAdmin(context.Background(), 999999999)
	s.Error(err)
	s.ErrorIs(err, auth.ErrInvalidCredentials)
}

func (s *AuthTestSuite) TestSetPriorityChannels_Error() {
	err := s.service.SetPriorityChannels(context.Background(), 999999, []string{"channel1"})
	s.Error(err)
}

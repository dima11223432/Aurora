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

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

package auth_test

import (
	"authService/internal/domain/models"
	"authService/internal/services/auth"
	"authService/internal/storage/postgres"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupStorage(t *testing.T) (*postgres.Storage, func()) {

	s, err := postgres.New("postgres://postgres:pass@localhost:5432/test_auth?sslmode=disable")
	require.NoError(t, err)

	teardown := func() {
		s.DB.Exec("TRUNCATE TABLE users, channels, apps RESTART IDENTITY CASCADE")
		s.DB.Close()
	}
	return s, teardown

}

func TestLoginUser(t *testing.T) {
	storage, teardown := setupStorage(t)
	defer teardown()

	service := auth.New(slog.Default(), storage, storage, storage, time.Duration(5*time.Minute))

	storage.DB.Exec("INSERT INTO apps (id, name, secret) VALUES (1, 'test', 'secret')")

	testCases := []struct {
		user    models.User
		WithErr bool
	}{
		{
			user: models.User{

				Telegram_id: 123456789,
				First_name:  "Dima",
				Last_name:   "Dmitriev",
				Username:    "dimadmitriev",
				Is_admin:    false,
			},
			WithErr: false,
		},

		{
			user: models.User{

				Telegram_id: 0,
				First_name:  "Dima",
				Last_name:   "Dmitriev",
				Username:    "dimadmitriev",
				Is_admin:    false,
			},
			WithErr: true,
		},

		{
			user: models.User{

				Telegram_id: 123456789,
				First_name:  "Dima",
				Last_name:   "",
				Username:    "dimadmitriev",
				Is_admin:    false,
			},
			WithErr: true,
		},
	}

	for _, tc := range testCases {
		token, err := service.Login(context.Background(), tc.user, 1)
		if tc.WithErr {
			assert.Error(t, err)
			assert.Empty(t, token)
		} else {
			assert.NoError(t, err)
			assert.NotEmpty(t, token)
		}
	}

}

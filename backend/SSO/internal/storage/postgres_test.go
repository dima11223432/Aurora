package storage_test

import (
	"authService/internal/domain/models"
	"authService/internal/storage/postgres"
	"context"
	"testing"

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

func TestStorage_SaveUser_and_User(t *testing.T) {

	s, teardown := setupStorage(t)
	defer teardown()

	ctx := context.Background()
	user := models.User{
		Telegram_id: 123456789,
		First_name:  "Dima",
		Last_name:   "Dmitriev",
		Username:    "dimadmitriev",
		Is_admin:    false,
	}

	id, err := s.SaveUser(ctx, user)
	user.ID = id

	assert.NoError(t, err)
	assert.NotZero(t, id)

	gottedUser, err := s.User(ctx, user.Telegram_id)
	assert.NoError(t, err)
	assert.Equal(t, gottedUser, user)
}

func TestStorage_SaveUser_and_User_empty(t *testing.T) {

	s, teardown := setupStorage(t)
	defer teardown()

	ctx := context.Background()
	user := models.User{
		Telegram_id: 123456789,
		First_name:  "any bullshit",
		Last_name:   "",
		Username:    "",
		Is_admin:    false,
	}

	id, err := s.SaveUser(ctx, user)
	user.ID = id

	assert.Error(t, err)
	assert.Zero(t, id)

	gottedUser, err := s.User(ctx, user.Telegram_id)
	assert.Error(t, err)
	assert.NotEqual(t, gottedUser, user)
}

package storage_test

import (
	"authService/internal/domain/models"
	"authService/internal/storage"
	"authService/internal/storage/postgres"
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PostgresTestSuite struct {
	suite.Suite
	storage *postgres.Storage
}

func (p *PostgresTestSuite) SetupTest() {

	s, err := postgres.New("postgres://postgres:pass@localhost:5432/test_auth?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	p.storage = s

	p.storage.DB.Exec("INSERT INTO apps (id, name, secret) VALUES (1, 'test', 'secret')")

}

func (p *PostgresTestSuite) TearDownTest() {

	p.storage.DB.Exec("TRUNCATE TABLE users, channels, apps RESTART IDENTITY CASCADE")
	p.storage.DB.Close()
}

func (p *PostgresTestSuite) TestStorage_SaveUser_and_User() {

	ctx := context.Background()
	user := models.User{
		Telegram_id: 123456789,
		First_name:  "Dima",
		Last_name:   "Dmitriev",
		Username:    "dimadmitriev",
		Is_admin:    false,
	}

	id, err := p.storage.SaveUser(ctx, user)
	user.ID = id

	p.NoError(err)
	p.NotZero(id)

	gottedUser, err := p.storage.User(ctx, user.Telegram_id)
	p.NoError(err)
	p.Equal(gottedUser, user)
}

func (p *PostgresTestSuite) TestStorage_SaveUser_and_User_empty() {

	ctx := context.Background()
	user := models.User{
		Telegram_id: 123456789,
		First_name:  "any bullshit",
		Last_name:   "",
		Username:    "",
		Is_admin:    false,
	}

	id, err := p.storage.SaveUser(ctx, user)
	user.ID = id

	p.Error(err)

	p.ErrorIs(err, storage.ErrEmptyUserValues)
	p.Zero(id)

	gottedUser, err := p.storage.User(ctx, user.Telegram_id)
	p.Error(err)
	p.ErrorIs(err, storage.ErrUserNotFound)
	p.NotEqual(gottedUser, user)
}

func (p *PostgresTestSuite) Test_storage_get_user_by_id() {

	user := models.User{
		Telegram_id: 123456789,
		First_name:  "Dima",
		Last_name:   "Dmitriev",
		Username:    "dimadmitriev",
		Is_admin:    false,
	}

	ctx := context.Background()
	id, err := p.storage.SaveUser(ctx, user)
	p.NoError(err)
	p.NotZero(id)
	user.ID = id

	gottedUser, err := p.storage.GetUserById(ctx, user.ID)
	p.NoError(err)
	p.Equal(gottedUser, user)
}

func (p *PostgresTestSuite) TestApp() {

	ctx := context.Background()
	p.storage.DB.Exec("INSERT INTO apps (id, name, secret) VALUES (1, 'test', 'secret')")

	app, err := p.storage.App(ctx, 1)
	p.NoError(err)
	p.Equal(app, models.App{ID: 1, Name: "test", Secret: "secret"})

	app, err = p.storage.App(ctx, 2)
	p.Error(err)
	p.ErrorIs(err, storage.ErrAppNotFound)
	p.Equal(app, models.App{})
}

func (p *PostgresTestSuite) TestSetPriorityChannels() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := models.User{
		Telegram_id: 123456789,
		First_name:  "Dima",
		Last_name:   "Dmitriev",
		Username:    "dimadmitriev",
		Is_admin:    false,
	}
	id, err := p.storage.SaveUser(ctx, user)
	p.NoError(err)
	p.NotZero(id)
	user.ID = id

	err = p.storage.SetPriorityChannels(ctx, user.ID, []string{"channel1", "channel2"})
	p.NoError(err)
}

func (p *PostgresTestSuite) TestSetPriorityChannelsEmptyChannels() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := models.User{
		Telegram_id: 123456789,
		First_name:  "Dima",
		Last_name:   "Dmitriev",
		Username:    "dimadmitriev",
		Is_admin:    false,
	}
	id, err := p.storage.SaveUser(ctx, user)
	p.NoError(err)
	p.NotZero(id)
	user.ID = id

	err = p.storage.SetPriorityChannels(ctx, user.ID, []string{})
	p.Error(err)
	p.ErrorIs(err, storage.ErrChannelsEmpty)
}

func (p *PostgresTestSuite) TestDeletePriorityChannels() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := models.User{
		Telegram_id: 123456789,
		First_name:  "Dima",
		Last_name:   "Dmitriev",
		Username:    "dimadmitriev",
		Is_admin:    false,
	}
	id, err := p.storage.SaveUser(ctx, user)
	p.NoError(err)
	p.NotZero(id)
	user.ID = id

	err = p.storage.SetPriorityChannels(ctx, user.ID, []string{"channel1", "channel2"})
	p.NoError(err)
	err = p.storage.DeletePriorityChannels(ctx, user.ID, []string{"channel1"})
	p.NoError(err)
}


func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}


func (p *PostgresTestSuite) TestStorage_IsAdmin() {
	ctx := context.Background()

	admin := models.User{
		Telegram_id: 1,
		First_name:  "A",
		Last_name:   "A",
		Username:    "admin",
		Is_admin:    true,
	}

	user := models.User{
		Telegram_id: 2,
		First_name:  "B",
		Last_name:   "B",
		Username:    "user",
		Is_admin:    false,
	}

	_, _ = p.storage.SaveUser(ctx, admin)
	_, _ = p.storage.SaveUser(ctx, user)

	isAdmin, err := p.storage.IsAdmin(ctx, 1)
	p.NoError(err)
	p.True(isAdmin)

	isAdmin, err = p.storage.IsAdmin(ctx, 2)
	p.NoError(err)
	p.False(isAdmin)
}

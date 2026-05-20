package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	pq "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type PostgresTestSuite struct {
	suite.Suite
	storage *Storage
}

func (p *PostgresTestSuite) SetupTest() {
	storagePass := os.Getenv("STORAGE_PASS")
	parsingServicePass := os.Getenv("PARSING_SERVICE_PASS")
	if storagePass == "" || parsingServicePass == "" {
		log.Fatal("STORAGE_PASS and PARSING_SERVICE_PASS environment variables are required")
	}
	s, err := New(storagePass, parsingServicePass)
	if err != nil {
		log.Fatal(err)
	}
	p.storage = s

}

func (p *PostgresTestSuite) TestGetAllParsingChannels() {
	ctx := context.Background()
	suffix := time.Now().UnixNano()
	channelName := fmt.Sprintf("test_channel_%d", suffix)

	_, err := p.storage.parserDB.ExecContext(ctx, `INSERT INTO default_channels (channel_username) VALUES ($1)`, channelName)
	p.NoError(err)

	p.T().Cleanup(func() {
		_, _ = p.storage.parserDB.ExecContext(context.Background(), `DELETE FROM default_channels WHERE channel_username = $1`, channelName)
	})

	channels, err := p.storage.GetAllDefaultParsingChannels(ctx)
	p.NoError(err)
	p.NotEmpty(channels)
}

func (p *PostgresTestSuite) TestDeleteParsingChannel() {
	ctx := context.Background()
	suffix := time.Now().UnixNano()
	channelName := fmt.Sprintf("test_delete_channel_%d", suffix)

	err := p.storage.AddNewParsingChannel(ctx, channelName)
	p.NoError(err)

	err = p.storage.DeleteDefaultParsingChannel(ctx, channelName)
	p.NoError(err)
}

func (p *PostgresTestSuite) TestGetPriorityChannelsByUserID() {
	ctx := context.Background()
	suffix := time.Now().UnixNano()
	telegramID := int64(9000000000 + suffix)
	channelName := fmt.Sprintf("test_priority_channel_%d", suffix)

	_, err := p.storage.db.ExecContext(ctx, `INSERT INTO users (telegram_id, first_name) VALUES ($1, $2)`, telegramID, "test")
	p.NoError(err)

	p.T().Cleanup(func() {
		_, _ = p.storage.db.ExecContext(context.Background(), `DELETE FROM users WHERE telegram_id = $1`, telegramID)
	})

	var userID int64
	err = p.storage.db.QueryRowContext(ctx, `SELECT user_id FROM users WHERE telegram_id = $1`, telegramID).Scan(&userID)
	p.NoError(err)

	_, err = p.storage.db.ExecContext(ctx, `INSERT INTO channels (user_id, channel_username) VALUES ($1, $2)`, userID, channelName)
	p.NoError(err)

	channels, err := p.storage.GetPriorityChannelsByUserID(ctx, userID)
	p.NoError(err)
	p.NotEmpty(channels)
}

func (p *PostgresTestSuite) TestGetAllCategoriesSuccessReturnsAllRecords() {
	ctx := context.Background()
	suffix := time.Now().UnixNano()
	cat1 := fmt.Sprintf("test_case1_cat_%d_a", suffix)
	cat2 := fmt.Sprintf("test_case1_cat_%d_b", suffix)

	_, err := p.storage.parserDB.ExecContext(
		ctx,
		`INSERT INTO channel_categories (name) VALUES ($1), ($2)`,
		cat1,
		cat2,
	)
	p.Require().NoError(err)

	p.T().Cleanup(func() {
		_, _ = p.storage.parserDB.ExecContext(
			context.Background(),
			`DELETE FROM channel_categories WHERE name IN ($1, $2)`,
			cat1,
			cat2,
		)
	})

	categories, err := p.storage.GetAllCategories(ctx)
	p.Require().NoError(err)
	p.Contains(categories, cat1)
	p.Contains(categories, cat2)
}

func (p *PostgresTestSuite) TestGetAllCategoriesScanString() {
	ctx := context.Background()
	suffix := time.Now().UnixNano()
	categoryName := fmt.Sprintf("test_case2_category_%d", suffix)

	_, err := p.storage.parserDB.ExecContext(
		ctx,
		`INSERT INTO channel_categories (name) VALUES ($1)`,
		categoryName,
	)
	p.Require().NoError(err)

	p.T().Cleanup(func() {
		_, _ = p.storage.parserDB.ExecContext(
			context.Background(),
			`DELETE FROM channel_categories WHERE name = $1`,
			categoryName,
		)
	})

	categories, err := p.storage.GetAllCategories(ctx)
	p.Require().NoError(err)
	p.Contains(categories, categoryName)
	for _, category := range categories {
		p.IsType("", category)
	}
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func TestGetDublicateErrorPostgresMatch(t *testing.T) {
	err := &pq.Error{Code: "23505"}

	isDuplicate := GetDublicateError(err)

	if !isDuplicate {
		t.Fatalf("expected true for postgres duplicate key code 23505")
	}
}

func TestGetDublicateErrorOtherError(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "regular error",
			err:  fmt.Errorf("some error"),
		},
		{
			name: "postgres other code",
			err:  &pq.Error{Code: "23503"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isDuplicate := GetDublicateError(tt.err)
			if isDuplicate {
				t.Fatalf("expected false for non-duplicate error")
			}
		})
	}
}

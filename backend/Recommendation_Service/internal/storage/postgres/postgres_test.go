package postgres

import (
	"context"
	"fmt"
	"log"
	"recommendationService/internal/config"
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
	cfg := config.MustLoadByPath("../../../config/local.yaml")
	s, err := New(cfg.StoragePass, cfg.ParsingServiceStoragePass)
	if err != nil {
		log.Fatal(err)
	}
	p.storage = s

}

func (p *PostgresTestSuite) TestGetAllParsingChannels() {
	ctx := context.Background()

	channels, err := p.storage.GetAllDefaultParsingChannels(ctx)
	p.NoError(err)
	p.NotEmpty(channels)
}

func (p *PostgresTestSuite) TestDeleteParsingChannel() {

	ctx := context.Background()
	channelName := "test_channel"

	err := p.storage.AddNewParsingChannel(ctx, channelName)
	p.NoError(err)

	err = p.storage.DeleteDefaultParsingChannel(ctx, channelName)
	p.NoError(err)
}

func (p *PostgresTestSuite) TestGetPriorityChannelsByUserID() {
	userID := int64(1)
	ctx := context.Background()
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
func (p *PostgresTestSuite) TestGetDefaultParsingChannelsByCategory() {
	ctx := context.Background()

	p.Run("success", func() {
		channels, err := p.storage.GetDefaultParsingChannelsByCategory(ctx, "news")

		p.NoError(err)
		p.NotEmpty(channels)

		for _, ch := range channels {
			p.Equal("news", ch.Category)
		}
	})

	p.Run("error", func() {
		channels, err := p.storage.GetDefaultParsingChannelsByCategory(nil, "news")

		p.Error(err)
		p.Empty(channels)
	})
}
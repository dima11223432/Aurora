package postgres

import (
	"context"
	"fmt"
	"log"
	"recommendationService/internal/config"
	"strings"
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

	channels, err := p.storage.GetAllParsingChannels(ctx)
	p.NoError(err)
	p.NotEmpty(channels)
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
		`INSERT INTO channel_categories (category_name) VALUES ($1), ($2)`,
		cat1,
		cat2,
	)
	p.Require().NoError(err)

	p.T().Cleanup(func() {
		_, _ = p.storage.parserDB.ExecContext(
			context.Background(),
			`DELETE FROM channel_categories WHERE category_name IN ($1, $2)`,
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
		`INSERT INTO channel_categories (category_name) VALUES ($1)`,
		categoryName,
	)
	p.Require().NoError(err)

	p.T().Cleanup(func() {
		_, _ = p.storage.parserDB.ExecContext(
			context.Background(),
			`DELETE FROM channel_categories WHERE category_name = $1`,
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

func TestGetAllCategoriesErrorReturnsNilAndOp(t *testing.T) {
	const op = "internal.storage.postgres.GetAllCategories"

	cfg := config.MustLoadByPath("../../../config/local.yaml")
	s, err := New(cfg.StoragePass, cfg.ParsingServiceStoragePass)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	if err := s.parserDB.Close(); err != nil {
		t.Fatalf("failed to close parser db: %v", err)
	}

	categories, err := s.GetAllCategories(context.Background())
	if categories != nil {
		t.Fatalf("expected nil categories on error, got: %#v", categories)
	}
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), op) {
		t.Fatalf("expected error to contain op %q, got: %v", op, err)
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

package postgres

import (
	"context"
	"log"
	"recommendationService/internal/config"
	"testing"

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
func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

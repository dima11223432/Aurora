package suite

import (
	"recomendationService/internal/config"
	postgresController "recomendationService/internal/storage/postgres"
	"testing"
)

type suite struct {
	*testing.T
	cfg                *config.Config
	PostgresController *postgresController.Storage
}

func NewSuite(t *testing.T) *suite {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local.yaml")

	postgresController, err := postgresController.New(cfg.StoragePass)
	if err != nil {
		t.Fatal(err)
	}
	return &suite{
		T:                  t,
		cfg:                cfg,
		PostgresController: postgresController,
	}
}

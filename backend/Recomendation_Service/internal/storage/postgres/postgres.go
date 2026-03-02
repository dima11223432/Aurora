package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"recomendationService/internal/domain/models"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

// New creates a new instance of the Storage
func New(storagePath string) (*Storage, error) {
	const op = "internal.storage.postgres.new"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{
		db: db,
	}, nil
}

func GetPriorityChannelsByUserID(ctx context.Context, userID int64) ([]models.PriorityChannel, error) {

	return nil, errors.New("Not implemented")
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"recommendationService/internal/domain/models"

	_ "github.com/lib/pq"
)

type Storage struct {
	db       *sql.DB
	parserDB *sql.DB
}

// New creates a new instance of the Storage
func New(storagePath string, parsingServicePath string) (*Storage, error) {
	const op = "internal.storage.postgres.new"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	parserDB, err := sql.Open("postgres", parsingServicePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{
		db:       db,
		parserDB: parserDB,
	}, nil
}

func (s *Storage) GetPriorityChannelsByUserID(ctx context.Context, userID int64) ([]models.PriorityChannel, error) {
	const op = "internal.storage.postgres.GetPriorityChannelsByUserID"

	q := `SELECT channel_username FROM channels WHERE user_id = $1`

	var channels []models.PriorityChannel

	query, err := s.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for query.Next() {
		var priorityChannel models.PriorityChannel

		if err := query.Scan(&priorityChannel.Channel); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		channels = append(channels, priorityChannel)
	}
	return channels, nil
}

func (s *Storage) GetAllParsingChannels(ctx context.Context) ([]string, error) {
	const op = "internal.storage.postgres.GetAllParsingChannels"

	q := `SELECT username FROM channels`

	channels := make([]string, 0)
	query, err := s.parserDB.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for query.Next() {
		var channel string
		if err := query.Scan(&channel); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		channels = append(channels, channel)
	}
	return channels, nil

}

func (s *Storage) AddNewParsingChannel(ctx context.Context, channel string) error {
	const op = "internal.storage.postgres.AddNewParsingChannel"

	q := `INSERT INTO channels (username) VALUES ($1)`
	_, err := s.parserDB.ExecContext(ctx, q, channel)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

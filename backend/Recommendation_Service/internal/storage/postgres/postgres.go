package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"recommendationService/internal/domain/models"
	"recommendationService/internal/storage"

	pq "github.com/lib/pq"
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

func (s *Storage) AddNewUserCustomParsingChannel(ctx context.Context, userID int64, channel string) error {
	const op = "internal.storage.postgres.AddNewUserCustomParsingChannel"
	q := `INSERT INTO user_custom_parsing_channels (user_id, channel_username) VALUES ($1, $2)`

	if _, err := s.parserDB.ExecContext(ctx, q, userID, channel); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.AddOnlyNewParsingChannel(ctx, channel)

	return nil
}

func (s *Storage) AddOnlyNewParsingChannel(ctx context.Context, channel string) error {
	const op = "internal.storage.postgres.AddOnlyNewParsingChannel"
	q := `INSERT INTO channels (username) VALUES ($1)`
	_, err := s.parserDB.ExecContext(ctx, q, channel)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, storage.ErrChannelExists)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) AddNewParsingChannel(ctx context.Context, channel string, category string) error {
	const op = "internal.storage.postgres.AddNewParsingChannel"

	q1 := `INSERT INTO channels (username) VALUES ($1)`
	q2 := `INSERT INTO channels_info (channel_id, category) VALUES ((SELECT id FROM channels WHERE username = $1), $2)`
	_, err := s.parserDB.ExecContext(ctx, q1, channel)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, storage.ErrChannelExists)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.parserDB.ExecContext(ctx, q2, channel, category)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteUserCustomParsingChannel(ctx context.Context, userID int64, channel string) error {
	const op = "internal.storage.postgres.DeleteUserCustomParsingChannel"
	q := `DELETE FROM user_custom_parsing_channels WHERE user_id = $1 AND channel_username = $2`
	_, err := s.parserDB.ExecContext(ctx, q, userID, channel)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	err = s.DeleteParsingChannel(ctx, channel)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteParsingChannel(ctx context.Context, channel string) error {
	const op = "internal.storage.postgres.DeleteParsingChannel"

	q := `DELETE FROM channels WHERE username = $1`
	_, err := s.parserDB.ExecContext(ctx, q, channel)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetParsingChannelsByCategory(ctx context.Context, category string) ([]string, error) {
	const op = "internal.storage.postgres.GetParsingChannelsByCategory"

	q := `SELECT c.username 
        FROM channels c
        INNER JOIN channels_info ci ON c.id = ci.channel_id
        WHERE ci.category = $1`
	query, err := s.parserDB.QueryContext(ctx, q, category)
	if err != nil {
		log.Println("DEHUIBHWEYUDGWIUYTEG")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	channelsUsernames := make([]string, 0)
	for query.Next() {
		var channel string
		if err := query.Scan(&channel); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		channelsUsernames = append(channelsUsernames, channel)
	}
	return channelsUsernames, nil
}

func (s *Storage) GetAllCategories(ctx context.Context) ([]string, error) {
	const op = "internal.storage.postgres.GetAllCategories"

	q := `SELECT name FROM channel_categories`
	query, err := s.parserDB.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	categories := make([]string, 0)
	for query.Next() {
		var category string
		if err := query.Scan(&category); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

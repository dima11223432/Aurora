package sqlite

import (
	"authService/internal/domain/models"
	"authService/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (s *Storage) SaveUser(ctx context.Context, user models.User) (int64, error) {
	const op = "storage.postgres.SaveUser"

	query := `
	INSERT INTO users (telegram_id, username, first_name, last_name, is_admin) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id
	`
	var userID int64

	err := s.db.QueryRowContext(ctx, query,
		user.Telegram_id,
		user.Username,
		user.First_name,
		user.Last_name,
		user.Is_admin,
	).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (s *Storage) User(ctx context.Context, telegram_id int64) (models.User, error) {
	const op = "storage.postgres.User"

	query := `
	SELECT id, telegram_id, username, first_name, last_name, is_admin 
	FROM users 
	WHERE telegram_id = $1
	`

	var user models.User
	err := s.db.QueryRowContext(ctx, query, telegram_id).Scan(
		&user.ID,
		&user.Telegram_id,
		&user.Username,
		&user.First_name,
		&user.Last_name,
		&user.Is_admin,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, telegram_id int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	query := `
	SELECT is_admin 
	FROM users 
	WHERE telegram_id = $1
	`

	var isAdmin bool

	err := s.db.QueryRowContext(ctx, query, telegram_id).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int64) (models.App, error) {
	const op = "storage.postgres.App"

	query := `
	SELECT id, name, secret 
	FROM apps 
	WHERE id = $1
	`

	var app models.App

	err := s.db.QueryRowContext(ctx, query, appID).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}

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

// NOTE:
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

func (s *Storage) SaveUser(ctx context.Context, email string, Passhash []byte, isAdmin bool) (int64, error) {
	const op = "storage.postgres.SaveUser"

	query := `
	INSERT INTO users (email, pass_hash, is_admin) VALUES ($1, $2, $3) RETURNING id
	`
	var userID int64

	err := s.db.QueryRowContext(ctx, query, email, Passhash, isAdmin).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	query := `
	SELECT id, email, pass_hash, is_admin FROM users WHERE email = $1
	`
	row := s.db.QueryRowContext(ctx, query, email)
	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.PassHash, &user.IsAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, sql.ErrNoRows)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	query := `
	SELECT is_admin FROM users WHERE id = $1
	`

	var IsAdmin bool

	err := s.db.QueryRowContext(ctx, query, userID).Scan(&IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return IsAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int64) (models.App, error) {
	const op = "storage.postgres.App"

	query := `
	SELECT id, name, secret FROM apps WHERE id = $1
	`

	var App models.App

	err := s.db.QueryRowContext(ctx, query, appID).Scan(&App.ID, &App.Name, &App.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("invalid appID")
		}
		return models.App{}, err
	}
	return App, nil
}

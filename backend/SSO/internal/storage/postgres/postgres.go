// Package postgres implements the PostgreSQL storage layer for the SSO service.
package postgres

import (
	"authService/internal/domain/models"
	"authService/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	pq "github.com/lib/pq"
)

var (
	emptyValue = 0
)

// Storage implements storage interfaces for users, apps, and channels using PostgreSQL.
type Storage struct {
	DB *sql.DB
}

// New creates a new instance of the Storage
func New(storagePath string) (*Storage, error) {
	const op = "internal.storage.postgres.new"

	DB, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := DB.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{
		DB: DB,
	}, nil
}

// SetPriorityChannels inserts priority channels for a user, skipping duplicates.
func (s *Storage) SetPriorityChannels(ctx context.Context, user_id int64, channels []string) error {
	const op = "storage.postgres.SetPriorityChannels"

	if len(channels) == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrChannelsEmpty)
	}

	query := `
	INSERT INTO channels (user_id, channel_username) VALUES
	`

	args := make([]interface{}, 0, len(channels)*2)

	for i, channel := range channels {
		query += fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)
		if i != len(channels)-1 {
			query += ","
		}
		args = append(args, user_id, channel)
	}
	query += " ON CONFLICT (user_id, channel_username) DO NOTHING"
	_, err := s.DB.ExecContext(ctx, query, args...)
	if err != nil {
		fmt.Println(err.Error())
		if isDuplicateError(err) {
			return fmt.Errorf("%s: %w", op, storage.ErrChannelExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// DeletePriorityChannels removes specified priority channels for a user.
func (s *Storage) DeletePriorityChannels(ctx context.Context, userID int64, channels []string) error {
	const op = "storage.postgres.DeletePriorityChannels"

	if len(channels) == 0 {
		return nil
	}

	query := `
	DELETE FROM channels 
	WHERE user_id = $1 AND channel_username = ANY($2);
	`

	_, err := s.DB.ExecContext(ctx, query, userID, pq.Array(channels))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// SaveUser inserts a new user into the database and returns the user ID.
func (s *Storage) SaveUser(ctx context.Context, user models.User) (int64, error) {
	const op = "storage.postgres.SaveUser"

	err := checkUserData(user)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	query := `
	INSERT INTO users (telegram_id, username, first_name, last_name, is_admin) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING user_id
	`
	var userID int64

	err = s.DB.QueryRowContext(ctx, query,
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

// GetUserById retrieves a user by their internal user ID.
func (s *Storage) GetUserById(ctx context.Context, user_id int64) (models.User, error) {
	const op = "storage.postgres.GetUserById"

	query := `
	SELECT user_id, telegram_id, username, first_name, last_name, is_admin 
	FROM users 
	WHERE user_id = $1
	`

	var user models.User
	err := s.DB.QueryRowContext(ctx, query, user_id).Scan(
		&user.ID, &user.Telegram_id, &user.Username, &user.First_name, &user.Last_name, &user.Is_admin,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

// User retrieves a user by their Telegram ID.
func (s *Storage) User(ctx context.Context, telegram_id int64) (models.User, error) {
	const op = "storage.postgres.User"

	query := `
	SELECT user_id, telegram_id, username, first_name, last_name, is_admin 
	FROM users 
	WHERE telegram_id = $1
	`

	var user models.User
	err := s.DB.QueryRowContext(ctx, query, telegram_id).Scan(
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

// IsAdmin checks the admin status of a user by their Telegram ID.
func (s *Storage) IsAdmin(ctx context.Context, telegram_id int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	query := `
	SELECT is_admin 
	FROM users 
	WHERE telegram_id = $1
	`

	var isAdmin bool

	err := s.DB.QueryRowContext(ctx, query, telegram_id).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isAdmin, nil
}

// App retrieves an application configuration by its ID.
func (s *Storage) App(ctx context.Context, appID int64) (models.App, error) {
	const op = "storage.postgres.App"

	query := `
	SELECT id, name, secret 
	FROM apps 
	WHERE id = $1
	`

	var app models.App

	err := s.DB.QueryRowContext(ctx, query, appID).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}
// isDuplicateError checks if a database error is a unique constraint violation.
func isDuplicateError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return strings.Contains(err.Error(), "duplicate key value")
}

// checkUserData validates that required user fields are non-empty.
func checkUserData(user models.User) error {
	if user.Telegram_id == int64(emptyValue) || user.First_name == "" || user.Last_name == "" || user.Username == "" {
		return storage.ErrEmptyUserValues
	}
	return nil
}

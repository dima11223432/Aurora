package auth

import (
	"authService/internal/domain/models"
	"authService/internal/lib/jwt"
	"authService/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type UserDataProvider struct {
	log                      *slog.Logger
	priorityChannelsProvider PriorityChannelsProvider
	TokenTTL                 time.Duration
}

type PriorityChannelsProvider interface {
	GetPriorityChannelsByUserID(ctx context.Context, userID int64)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

// New returns a new instance of Auth service
func New(
	log *slog.Logger,
	priorityChannelsProvider PriorityChannelsProvider,

	tokenTTL time.Duration,
) *UserDataProvider {
	return &UserDataProvider{
		log:                      log,
		priorityChannelsProvider: priorityChannelsProvider,
		TokenTTL:                 tokenTTL,
	}
}

func GetUserPriorityChanneld(ctx context.Context, userID int64) ([]string, error) {
	return nil, errors.New("Not implemented")
}

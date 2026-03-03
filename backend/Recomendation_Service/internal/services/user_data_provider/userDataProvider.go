package auth

import (
	"context"
	"errors"
	"log/slog"
	"recomendationService/internal/domain/models"
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

func (u *UserDataProvider) GetUserPriorityChanneld(ctx context.Context, userID int64) ([]models.PriorityChannel, error) {
	return nil, errors.New("Not implemented")
}

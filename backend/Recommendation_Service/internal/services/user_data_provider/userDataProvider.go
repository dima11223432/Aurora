package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"recommendationService/internal/domain/models"
	"time"
)

type UserDataProvider struct {
	log                      *slog.Logger
	priorityChannelsProvider PriorityChannelsProvider
	TokenTTL                 time.Duration
}

type PriorityChannelsProvider interface {
	GetPriorityChannelsByUserID(ctx context.Context, userID int64) ([]models.PriorityChannel, error)
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

func (u *UserDataProvider) GetUserPriorityChannels(ctx context.Context, userID int64) ([]models.PriorityChannel, error) {
	const op = "internal.services.user_data_provider.userDataProvider.go.GetUserPriorityChannels"

	channels, err := u.priorityChannelsProvider.GetPriorityChannelsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return channels, nil
}

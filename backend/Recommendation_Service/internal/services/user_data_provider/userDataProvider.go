package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"recommendationService/internal/domain/models"
	"time"
)

type UserDataProvider struct {
	log                      *slog.Logger
	priorityChannelsProvider PriorityChannelsProvider
	priorityNewsProvider     PriorityNewsProvider
	TokenTTL                 time.Duration
}

type PriorityChannelsProvider interface {
	GetPriorityChannelsByUserID(ctx context.Context, userID int64) ([]models.PriorityChannel, error)
}

type PriorityNewsProvider interface {
	GetPostsByChannels(ctx context.Context, channels []string, userID int64, cursor *models.Cursor, limit int64) ([]models.Post, *models.Cursor, error)
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
	priorityNewsProvider PriorityNewsProvider,
	tokenTTL time.Duration,
) *UserDataProvider {
	return &UserDataProvider{
		log:                      log,
		priorityChannelsProvider: priorityChannelsProvider,
		priorityNewsProvider:     priorityNewsProvider,
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

func (u *UserDataProvider) GetRecommendatedPosts(ctx context.Context, userID int64) ([]models.Post, error) {
	const op = "internal.services.user_data_provider.userDataProvider.go.GetUserPriorityNews"

	channels, err := u.GetUserPriorityChannels(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(channels) == 0 {
		return []models.Post{}, nil
	}

	channelNames := make([]string, 0, len(channels))
	for _, ch := range channels {
		channelNames = append(channelNames, ch.Channel)
	}

	posts, nextCursor, err := u.priorityNewsProvider.GetPostsByChannels(ctx, channelNames, userID, nil, 10)
	log.Print(nextCursor)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nil
}

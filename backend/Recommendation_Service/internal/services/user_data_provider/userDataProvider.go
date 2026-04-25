package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"recommendationService/internal/domain/models"
	"recommendationService/internal/storage"
	"time"
)

type UserDataProvider struct {
	log                      *slog.Logger
	priorityChannelsProvider PriorityChannelsProvider
	parsingChannelsProvider  ParsingChannelsProvider
	priorityNewsProvider     PriorityNewsProvider
	TokenTTL                 time.Duration
}

type PriorityChannelsProvider interface {
	GetPriorityChannelsByUserID(ctx context.Context, userID int64) ([]models.PriorityChannel, error)
}

type ParsingChannelsProvider interface {
	GetAllParsingChannels(ctx context.Context) ([]string, error)
	AddNewParsingChannel(ctx context.Context, channel string) error
	DeleteParsingChannel(ctx context.Context, channel string) error
	GetAllCategories(ctx context.Context) ([]string, error)
	GetParsingChannelsByCategory(ctx context.Context, category string) ([]string, error)
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
	parsingChannelsProvider ParsingChannelsProvider,
	priorityNewsProvider PriorityNewsProvider,
	tokenTTL time.Duration,
) *UserDataProvider {
	return &UserDataProvider{
		log:                      log,
		priorityChannelsProvider: priorityChannelsProvider,
		priorityNewsProvider:     priorityNewsProvider,
		parsingChannelsProvider:  parsingChannelsProvider,
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

func (u *UserDataProvider) GetRecommendatedPosts(ctx context.Context, userID int64, cursor *models.Cursor) ([]models.Post, *models.Cursor, error) {
	const op = "internal.services.user_data_provider.userDataProvider.go.GetUserPriorityNews"

	channels, err := u.GetUserPriorityChannels(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(channels) == 0 {
		return []models.Post{}, nil, nil
	}

	channelNames := make([]string, 0, len(channels))
	for _, ch := range channels {
		channelNames = append(channelNames, ch.Channel)
	}

	posts, nextCursor, err := u.priorityNewsProvider.GetPostsByChannels(ctx, channelNames, userID, cursor, 4)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nextCursor, nil
}

func (u *UserDataProvider) GetAllParsingChannels(ctx context.Context) ([]string, error) {
	const op = "internal.services.user_data_provider.userDataProvider.go.GetAllParsingChannels"
	channels, err := u.parsingChannelsProvider.GetAllParsingChannels(ctx)
	if err != nil {
		u.log.Error("failed to get all parsing channels",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return channels, nil
}

func (u *UserDataProvider) AddNewParsingChannel(ctx context.Context, channel string) error {
	const op = "internal.services.user_data_provider.userDataProvider.go.AddNewParsingChannel"
	err := u.parsingChannelsProvider.AddNewParsingChannel(ctx, channel)
	if err != nil {
		if errors.Is(err, storage.ErrChannelExists) {
			u.log.Error("parsing channel already exists",
				slog.String("op", op),
				slog.String("channel", channel),
				slog.Any("err", err),
			)
			return fmt.Errorf("%s: %w", op, storage.ErrChannelExists)
		}
		u.log.Error("failed to add new parsing channel",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (u *UserDataProvider) DeleteParsingChannel(ctx context.Context, channel string) error {
	const op = "internal.services.user_data_provider.userDataProvider.go.DeleteParsingChannel"
	err := u.parsingChannelsProvider.DeleteParsingChannel(ctx, channel)
	if err != nil {
		u.log.Error("failed to delete parsing channel",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (u *UserDataProvider) GetParsingChannelsWithCategories(ctx context.Context) (map[string][]string, error) {
	const op = "internal.services.user_data_provider.userDataProvider.go.GetParsingChannelsWithCategories"

	categoriesWithChannels := make(map[string][]string, 0)

	categories, err := u.parsingChannelsProvider.GetAllCategories(ctx)
	if err != nil {
		u.log.Error("failed to get parsing channels with categories",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for _, category := range categories {
		channels, err := u.parsingChannelsProvider.GetParsingChannelsByCategory(ctx, category)
		if err != nil {
			u.log.Error("failed to get parsing channels with categories",
				slog.String("op", op),
				slog.Any("err", err),
			)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		categoriesWithChannels[category] = channels
	}
	return categoriesWithChannels, nil
}

func (u *UserDataProvider) GetParsingChannelsByCategory(ctx context.Context, category string) ([]string, error) {
	const op = "internal.services.user_data_provider.userDataProvider.go.GetParsingChannelsByCategory"
	channels, err := u.parsingChannelsProvider.GetParsingChannelsByCategory(ctx, category)
	if err != nil {
		u.log.Error("failed to get parsing channels by category",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return channels, nil
}

func (u *UserDataProvider) GetAllCategories(ctx context.Context) ([]string, error) {
	const op = "internal.services.user_data_provider.userDataProvider.go.GetAllCategories"
	categories, err := u.parsingChannelsProvider.GetAllCategories(ctx)
	if err != nil {
		u.log.Error("failed to get all categories",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return categories, nil
}

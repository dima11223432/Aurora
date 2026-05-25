// Package auth (user_data_provider) implements business logic for user data management,
// including priority channels, parsing channels, and recommended news posts.
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

// UserDataProvider orchestrates user-related data operations by delegating
// to priority channel, parsing channel, and priority news providers.
type UserDataProvider struct {
	log                      *slog.Logger
	priorityChannelsProvider PriorityChannelsProvider
	parsingChannelsProvider  ParsingChannelsProvider
	priorityNewsProvider     PriorityNewsProvider
	TokenTTL                 time.Duration
}

// PriorityChannelsProvider defines storage operations for priority channels.
type PriorityChannelsProvider interface {
	GetPriorityChannelsByUserID(ctx context.Context, userID int64) ([]models.PriorityChannel, error)
}

// ParsingChannelsProvider defines storage operations for parsing channels.
type ParsingChannelsProvider interface {
	GetAllDefaultParsingChannels(ctx context.Context) ([]string, error)
	AddNewDefaultParsingChannel(ctx context.Context, channel string, category string) error
	DeleteDefaultParsingChannel(ctx context.Context, channel string) error
	DeleteUserCustomParsingChannel(ctx context.Context, userID int64, channel string) error
	DeleteParsingChannel(ctx context.Context, channel string) error
	GetAllUserCustomParsingChannels(ctx context.Context, userID int64) ([]string, error)
	GetAllCategories(ctx context.Context) ([]string, error)
	AddNewUserCustomParsingChannel(ctx context.Context, userID int64, channel string) error
	GetDefaultParsingChannelsByCategory(ctx context.Context, category string) ([]string, error)
	AddNewParsingChannel(ctx context.Context, channel string) error
	SetChannelCategory(ctx context.Context, channel string, category string) error
}

// PriorityNewsProvider defines storage operations for retrieving news posts by channels.
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

// GetUserPriorityChannels returns priority channels for a given user.
func (u *UserDataProvider) GetUserPriorityChannels(ctx context.Context, userID int64) ([]models.PriorityChannel, error) {
	const op = "internal.services.user_data_provider.userDataProvider.go.GetUserPriorityChannels"

	channels, err := u.priorityChannelsProvider.GetPriorityChannelsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return channels, nil
}

// GetRecommendatedPosts fetches recommended posts for a user based on their priority channels.
// Supports cursor-based pagination.
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

	posts, nextCursor, err := u.priorityNewsProvider.GetPostsByChannels(ctx, channelNames, userID, cursor, 15)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nextCursor, nil
}

// GetAllDefaultParsingChannels returns all system-wide default parsing channels.
func (u *UserDataProvider) GetAllDefaultParsingChannels(ctx context.Context) ([]string, error) {
	const op = "internal.services.user_data_provider.userDataProvider.go.GetAllParsingChannels"
	channels, err := u.parsingChannelsProvider.GetAllDefaultParsingChannels(ctx)
	if err != nil {
		u.log.Error("failed to get all parsing channels",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return channels, nil
}

// AddNewUserCustomParsingChannel adds a custom parsing channel for a user.
func (u *UserDataProvider) AddNewUserCustomParsingChannel(ctx context.Context, userID int64, channel string) error {
	const op = "internal.services.user_data_provider.userDataProvider.go.AddNewUserCustomParsingChannel"
	if err := u.parsingChannelsProvider.AddNewUserCustomParsingChannel(ctx, userID, channel); err != nil {
		u.log.Error("failed to add new user custom parsing channel",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

// AddNewDefaultParsingChannel adds a new system-wide default parsing channel with a category.
func (u *UserDataProvider) AddNewDefaultParsingChannel(ctx context.Context, channel string, category string) error {
	const op = "internal.services.user_data_provider.userDataProvider.go.AddNewParsingChannel"
	err := u.parsingChannelsProvider.AddNewDefaultParsingChannel(ctx, channel, category)
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
	err = u.parsingChannelsProvider.AddNewParsingChannel(ctx, channel)
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
	err = u.parsingChannelsProvider.SetChannelCategory(ctx, channel, category)
	if err != nil {
		u.log.Error("failed to set channel category",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// DeleteUserCustomParsingChannel removes a custom parsing channel for a user.
func (u *UserDataProvider) DeleteUserCustomParsingChannel(ctx context.Context, userID int64, channel string) error {
	const op = "internal.services.user_data_provider.userDataProvider.go.DeleteUserCustomParsingChannel"
	if err := u.parsingChannelsProvider.DeleteUserCustomParsingChannel(ctx, userID, channel); err != nil {
		u.log.Error("failed to delete user custom parsing channel",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := u.parsingChannelsProvider.DeleteParsingChannel(ctx, channel); err != nil {
		u.log.Error("failed to delete parsing channel",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// DeleteDefaultParsingChannel removes a system-wide default parsing channel.
func (u *UserDataProvider) DeleteDefaultParsingChannel(ctx context.Context, channel string) error {
	const op = "internal.services.user_data_provider.userDataProvider.go.DeleteParsingChannel"
	err := u.parsingChannelsProvider.DeleteDefaultParsingChannel(ctx, channel)
	if err != nil {
		u.log.Error("failed to delete default parsing channel",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}
	err = u.parsingChannelsProvider.DeleteParsingChannel(ctx, channel)
	if err != nil {
		u.log.Error("Failed to delete parsing channel",
			slog.String("op", op),
			slog.Any("err", err))
	}
	return nil
}

// GetAllUserCustomParsingChannels returns all custom parsing channels for a user.
func (u *UserDataProvider) GetAllUserCustomParsingChannels(ctx context.Context, userID int64) ([]string, error) {
	const op = "internal.services.user_dats_provider.userDataProvider.go.GetAllDefaultParsingChannels"
	channels, err := u.parsingChannelsProvider.GetAllUserCustomParsingChannels(ctx, userID)
	if err != nil {
		u.log.Error("failed to get all user custom parsing channels",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return channels, nil
}

// GetDefaultParsingChannelsWithCategories returns default parsing channels grouped by category.
func (u *UserDataProvider) GetDefaultParsingChannelsWithCategories(ctx context.Context) (map[string][]string, error) {
	const op = "internal.services.user_data_provider.GetDefaultParsingChannelsWithCategories"

	categories, err := u.parsingChannelsProvider.GetAllCategories(ctx)
	if err != nil {
		u.log.Error("step 1: failed to get all categories", slog.String("op", op), slog.Any("err", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	categoriesWithChannels := make(map[string][]string)

	for _, category := range categories {
		channels, err := u.parsingChannelsProvider.GetDefaultParsingChannelsByCategory(ctx, category)
		if err != nil {
			u.log.Error("step 2: failed to get channels for category",
				slog.String("category", category),
				slog.Any("err", err),
			)
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if channels == nil {
			channels = []string{}
		}
		categoriesWithChannels[category] = channels
	}

	return categoriesWithChannels, nil
}

// GetDefaultParsingChannelsByCategory returns parsing channels for a given category.
func (u *UserDataProvider) GetDefaultParsingChannelsByCategory(ctx context.Context, category string) ([]string, error) {
	const op = "internal.services.user_data_provider.userDataProvider.go.GetParsingChannelsByCategory"
	channels, err := u.parsingChannelsProvider.GetDefaultParsingChannelsByCategory(ctx, category)
	if err != nil {
		u.log.Error("failed to get parsing channels by category",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return channels, nil
}

// GetAllCategories returns all available channel categories.
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

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

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	TokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, user models.User) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, telegram_id int64) (models.User, error)
	IsAdmin(ctx context.Context, telegram_id int64) (bool, error)
	SetPriorityChannels(ctx context.Context, user_id int64, channels []string) error
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

// New returns a new instance of Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userProvider: userProvider,
		userSaver:    userSaver,
		appProvider:  appProvider,
		TokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, user models.User, appID int) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("telegram_id", user.Telegram_id),
		slog.String("username", user.Username),
		slog.String("first_name", user.First_name),
	)
	log.Info("attempting to login user")

	dbUser, err := a.userProvider.User(ctx, user.Telegram_id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found, registering new user")

			id, regErr := a.RegisterNewUser(ctx, user)
			if regErr != nil {
				return "", fmt.Errorf("%s: %w", op, regErr)
			}

			dbUser, err = a.userProvider.User(ctx, user.Telegram_id)
			if err != nil {
				return "", fmt.Errorf("%s: %w", op, err)
			}

			log.Info("user registered successfully",
				slog.Int64("user_id", id),
				slog.Int64("telegram_id", dbUser.Telegram_id),
			)
		} else {
			log.Error("failed to get user", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	app, err := a.appProvider.App(ctx, int64(appID))
	if err != nil {
		log.Error("failed to get app", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(dbUser, app, a.TokenTTL)
	if err != nil {
		log.Error("failed to create token", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully",
		slog.Int64("user_id", dbUser.ID),
		slog.Int64("telegram_id", dbUser.Telegram_id),
	)
	return token, nil
}

//func SetPriorityChannels
func (a *Auth) SetPriorityChannels(ctx context.Context, user_id int64, channels []string) (int32, error) {
	const op = "auth.SetPriorityChannels"

	err := a.UserProvider.SetPriorityChannels(ctx context.Context, user_id int64, channels []string)
	if err != nil {
		return 400, fmt.Errorf("%s: %w", op, err)
	}
	return 200, nil

}



func (a *Auth) RegisterNewUser(ctx context.Context, user models.User) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("telegram_id", user.Telegram_id),
		slog.String("username", user.Username),
		slog.String("first_name", user.First_name),
	)

	log.Info("attempting to register new user")

	existingUser, err := a.userProvider.User(ctx, user.Telegram_id)
	if err == nil {
		log.Warn("user already exists", slog.Int64("existing_id", existingUser.ID))
		return existingUser.ID, fmt.Errorf("%s: %w", op, ErrUserExists)
	}
	if !errors.Is(err, storage.ErrUserNotFound) {
		log.Error("failed to check user existence", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, user)
	if err != nil {
		log.Error("failed to save user", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered successfully",
		slog.Int64("user_id", id),
		slog.Int64("telegram_id", user.Telegram_id),
	)
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, telegram_id int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("telegram_id", telegram_id),
	)
	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, telegram_id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found")
			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to check admin status", slog.String("error", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))
	return isAdmin, nil
}

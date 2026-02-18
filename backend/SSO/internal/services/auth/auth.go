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
	)
	log.Info("attempting to login User")

	user, err := a.userProvider.User(ctx, user.Telegram_id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			id, err := a.RegisterNewUser(ctx, user)
			if err != nil {
				return "", fmt.Errorf("%s, %w", op, err)
			}
			user, err = a.userProvider.User(ctx, user.Telegram_id)
			if err != nil {
				return "", fmt.Errorf("%s, %w", op, err)
			}

			slog.Info("user registered successfully", slog.Int64("id", id))

		} else {
			return "", fmt.Errorf("%s, %w", op, err)
		}
	}

	app, err := a.appProvider.App(ctx, int64(appID))
	if err != nil {
		return "", fmt.Errorf("%s, %w", op, err)
	}
	log.Info("user logged in successfully")
	token, err := jwt.NewToken(user, app, a.TokenTTL)
	if err != nil {
		a.log.Error("failed to create token")
		return "", fmt.Errorf("%s, %w", op, err)
	}
	return token, nil
}
func (a *Auth) RegisterNewUser(ctx context.Context, user models.User) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("email", user.Telegram_id),
	)

	log.Info("attempting to register new user")

	id, err := a.userSaver.SaveUser(ctx, user)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists")
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered successfully", slog.Int64("user_id", id))
	return id, nil
}
func (a *Auth) IsAdmin(ctx context.Context, telegram_id int64) (bool, error) {
	const op = "auth.IsAdmin"
	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", telegram_id),
	)
	log.Info("checking id user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, telegram_id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found")
			return false, fmt.Errorf("%s, %w", op, ErrInvalidAppID)
		}
		return false, fmt.Errorf("%s, %w", op, err)
	}
	slog.Info("checked id user is admin", slog.Bool("isAdmin", isAdmin))
	return isAdmin, nil
}

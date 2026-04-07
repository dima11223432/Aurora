package services

import (
	custom_errors "API_Service/internal/custom_errors"
	"context"
	"fmt"
	"log/slog"
	"time"

	// ssov1 "github.com/dima11223432/protos/gen/go/sso"
	ssov1 "github.com/dima11223432/Aurora_SSO_Protos/api/gen/v1"
	// recv1 "github.com/dima11223432/recommendationService_protos/api/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthInterceptor interface {
	SetAuthInterceptor() grpc.UnaryServerInterceptor
	GetUserIdFromContext(ctx context.Context) (int64, error)
}
type AuthService struct {
	log             *slog.Logger
	AuthClient      ssov1.AuthServiceClient
	AuthInterceptor AuthInterceptor
}

func NewAuthService(log *slog.Logger, authClient ssov1.AuthServiceClient, authinterceptor AuthInterceptor) *AuthService {
	return &AuthService{
		log:             log,
		AuthClient:      authClient,
		AuthInterceptor: authinterceptor,
	}
}

func (a *AuthService) DeletePriorityChannels(ctx context.Context, channels []string) error {
	const op = "ApiService.internal.services.AuthService.DeletePriorityChannels"

	userID, err := a.AuthInterceptor.GetUserIdFromContext(ctx)
	if err != nil {
		a.log.Error("invalid user id in context", slog.String("op", op), slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = a.AuthClient.DeletePriorityChannels(ctx, &ssov1.DeletePriorityChannelsRequest{
		UserId:            userID,
		ChannelsUsernames: channels,
	})
	if err != nil {
		a.log.Error("failed to delete user priority channels", slog.String("op", op), slog.Any("err", err))
	}
	return nil
}

func (a *AuthService) SetPriorityChannels(ctx context.Context, channels []string) error {
	const op = "services.AuthService.SetPriorityChannels"

	userID, err := a.AuthInterceptor.GetUserIdFromContext(ctx)
	if err != nil {
		a.log.Error("failed to get user id from context",
			slog.String("op", op),
			slog.Any("err", err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = a.AuthClient.SetPriorityChannels(ctx, &ssov1.SetPriorityChannelsRequest{
		UserId:            userID,
		ChannelsUsernames: channels,
	})

	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return custom_errors.ErrChannelExists
		}

		a.log.Error("failed to set priority channels",
			slog.String("op", op),
			slog.Int64("user_id", userID),
			slog.Any("err", err),
		)

		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *AuthService) Login(
	ctx context.Context,
	telegram_id int64,
	username string,
	firstName string,
	lastName string,
	appId int64,
) (string, error) {
	resp, err := a.AuthClient.Login(ctx, &ssov1.LoginRequest{
		TelegramId: telegram_id,
		Username:   username,
		AppId:      appId,
		FirstName:  firstName,
		LastName:   lastName,
		IsAdmin:    false,
	})

	if err != nil {
		a.log.Error("failed to login user", slog.String("error", err.Error()))
		return "", err
	}

	return resp.Token, nil
}

func (a *AuthService) IsAdmin(
	ctx context.Context,
	telegram_id int64,
) (bool, error) {
	resp, err := a.AuthClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
		TelegramId: telegram_id,
	})
	if err != nil {
		a.log.Error("failed to check admin status", slog.String("error", err.Error()))
		return false, err
	}

	return resp.IsAdmin, nil
}

package services

import (
	"context"

	// ssov1 "github.com/dima11223432/protos/gen/go/sso"
	ssov1 "github.com/dima11223432/Aurora_SSO_Protos/api/gen/v1"
	"google.golang.org/grpc"
)

type AuthInterceptor interface {
	SetAuthInterceptor() grpc.UnaryServerInterceptor
	GetUserIdFromContext(ctx context.Context) (int64, error)
}
type AuthService struct {
	AuthClient      ssov1.AuthServiceClient
	AuthInterceptor AuthInterceptor
}

func NewAuthService(authClient ssov1.AuthServiceClient, authinterceptor AuthInterceptor) *AuthService {
	return &AuthService{
		AuthClient:      authClient,
		AuthInterceptor: authinterceptor,
	}
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
		return false, err
	}

	return resp.IsAdmin, nil
}

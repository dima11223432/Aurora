package services

import (
	"context"

	ssov1 "github.com/dima11223432/protos/gen/go/sso"
)

type AuthService struct {
	AuthClient ssov1.AuthClient
}

func NewAuthService(authClient ssov1.AuthClient) *AuthService {
	return &AuthService{
		AuthClient: authClient,
	}
}

func (a *AuthService) Register(
	ctx context.Context,
	email string,
	password string,
	is_admin bool,
) (int64, error) {
	resp, err := a.AuthClient.Register(ctx, &ssov1.RegisterRequest{Email: email, Password: password, IsAdmin: is_admin})
	if err != nil {
		return 0, err
	}

	return resp.UserId, nil
}

func (a *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	appId int32,
) (string, error) {
	resp, err := a.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appId,
	})

	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

func (a *AuthService) IsAdmin(
	ctx context.Context,
	userId int64,
) (bool, error) {
	resp, err := a.AuthClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: userId,
	})
	if err != nil {
		return false, err
	}

	return resp.IsAdmin, nil
}

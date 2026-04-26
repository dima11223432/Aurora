package authinterceptor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type CtxUserIdKey struct{}
type CtxIsAdminKey struct{}

type AuthInterceptor struct {
	AuthConfig AuthConfig
}

type Claims struct {
	ID      int64 `json:"id"`
	IsAdmin bool  `json:"is_admin"`
	jwt.RegisteredClaims
}

type AuthConfig struct {
	JwtSecret    string
	PublicRoutes []string
}

func NewAuthInterceptor(authConfig AuthConfig) *AuthInterceptor {
	return &AuthInterceptor{
		AuthConfig: authConfig,
	}
}

func (a *AuthInterceptor) parseToken(tokenStr, secret string) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, err := parser.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token or claims type")
	}

	return claims, nil
}

func (a *AuthInterceptor) SetAuthInterceptor() grpc.UnaryServerInterceptor {
	publicMethods := make(map[string]bool)
	for _, method := range a.AuthConfig.PublicRoutes {
		publicMethods[method] = true
	}
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing auth header")
		}

		tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")

		claims, err := a.parseToken(tokenStr, a.AuthConfig.JwtSecret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "can't parse jwt token")
		}

		ctx = context.WithValue(ctx, CtxUserIdKey{}, claims.ID)
		ctx = context.WithValue(ctx, CtxIsAdminKey{}, claims.IsAdmin)

		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) GetUserIdFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(CtxUserIdKey{}).(int64)
	if !ok {
		return 0, errors.New("user id not found in context")
	}
	return userID, nil
}

func (a *AuthInterceptor) GetIsAdminFromContext(ctx context.Context) (bool, error) {
	isAdmin, ok := ctx.Value(CtxIsAdminKey{}).(bool)
	if !ok {
		return false, errors.New("is admin not found in context")
	}
	return isAdmin, nil
}

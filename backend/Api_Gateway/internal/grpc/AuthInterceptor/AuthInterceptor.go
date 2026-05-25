// Package authinterceptor provides gRPC unary interceptor for JWT-based authentication.
// It extracts user identity from incoming requests and injects it into the context.
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

// CtxUserIdKey is the context key for storing the authenticated user ID.
type CtxUserIdKey struct{}

// CtxIsAdminKey is the context key for storing the admin status of the authenticated user.
type CtxIsAdminKey struct{}

// AuthInterceptor provides JWT authentication for gRPC requests.
// It intercepts incoming calls, validates JWT tokens, and injects
// user ID and admin status into the context.
type AuthInterceptor struct {
	AuthConfig AuthConfig
}

// Claims represents the JWT token claims used for authentication.
type Claims struct {
	ID      int64 `json:"id"`
	IsAdmin bool  `json:"is_admin"`
	jwt.RegisteredClaims
}

// AuthConfig holds configuration for the authentication interceptor.
type AuthConfig struct {
	JwtSecret    string
	PublicRoutes []string
}

// NewAuthInterceptor creates a new AuthInterceptor with the given configuration.
func NewAuthInterceptor(authConfig AuthConfig) *AuthInterceptor {
	return &AuthInterceptor{
		AuthConfig: authConfig,
	}
}

// parseToken parses and validates a JWT token string using the given secret.
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

// SetAuthInterceptor returns a gRPC UnaryServerInterceptor that validates JWT tokens
// for non-public routes and injects user context.
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

// GetUserIdFromContext extracts the authenticated user ID from the context.
func (a *AuthInterceptor) GetUserIdFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(CtxUserIdKey{}).(int64)
	if !ok {
		return 0, errors.New("user id not found in context")
	}
	return userID, nil
}

// GetIsAdminFromContext extracts the admin status of the authenticated user from context.
func (a *AuthInterceptor) GetIsAdminFromContext(ctx context.Context) (bool, error) {
	isAdmin, ok := ctx.Value(CtxIsAdminKey{}).(bool)
	if !ok {
		return false, errors.New("is admin not found in context")
	}
	return isAdmin, nil
}

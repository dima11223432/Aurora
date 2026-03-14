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

type AnasyledDataProvider interface {
	SetValue(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error
	GetValue(ctx context.Context, key string) (interface{}, error)
}

type RedisService struct {
	log *slog.Logger
	AnasyledDataProvider
	TokenTTL time.Duration
}

func NewRedisService(log *slog.Logger, tokenTTL time.Duration) *RedisService {
	return &RedisService{
		log:      log,
		TokenTTL: tokenTTL,
	}
}

func (r *RedisService) SetAnalysedData(ctx context.Context, dataTitle string, analysedData interface{}) error {
	const op = "Cache_Service.internal.services.auth.SetAnalysedData"

	err := r.AnasyledDataProvider.SetValue(ctx, dataTitle, analysedData, r.TokenTTL)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}

// TODO: ревлизовать метод, вызывающий GetValue из AnasyledDataProvider
func (r *RedisService) GetAnalyseData(ctx context.Context, dataTitle string) (interface{}, error) {
	return nil, nil
}

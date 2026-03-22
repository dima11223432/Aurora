package analyseDataProvider

import (
	"CacheService/internal/domain/models"
	"context"
	"fmt"
	"log/slog"
	"time"
)

type AnasyledDataProvider interface {
	SetValue(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error
	GetValue(ctx context.Context, key string) (interface{}, error)
}

type RedisService struct {
	log      *slog.Logger
	provider AnasyledDataProvider
	TokenTTL time.Duration
}

func NewRedisService(log *slog.Logger, analyseDataProvider AnasyledDataProvider, tokenTTL time.Duration) *RedisService {
	return &RedisService{
		log:      log,
		TokenTTL: tokenTTL,
		provider: analyseDataProvider,
	}
}

func (r *RedisService) SetAnalysedData(ctx context.Context, dataTitle string, analysedData interface{}) error {
	const op = "Cache_Service.internal.services.auth.SetAnalysedData"

	err := r.provider.SetValue(ctx, dataTitle, analysedData, r.TokenTTL)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	r.log.Info("SetAnalysedData", slog.String("dataTitle", dataTitle))
	return nil

}

func (r *RedisService) GetAnalysedData(ctx context.Context, dataTitle string) (interface{}, error) {
	data, err := r.provider.GetValue(ctx, dataTitle)
	if err != nil {
		r.log.Error("Failed to get analysedData", slog.Any("error", err))
		return nil, err
	}
	r.log.Info("GetAnalysedData", slog.String("dataTitle", dataTitle))
	return data, err
}

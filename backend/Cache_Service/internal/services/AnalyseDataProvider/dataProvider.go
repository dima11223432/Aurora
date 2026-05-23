// Package analyseDataProvider provides business logic for storing and retrieving
// analysed financial data using a Redis-backed cache.
package analyseDataProvider

import (
	"CacheService/internal/domain/models"
	"context"
	"fmt"
	"log/slog"
	"time"
)

// AnasyledDataProvider defines cache operations for analysed data.
type AnasyledDataProvider interface {
	SetValue(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error
	GetValue(ctx context.Context, key string) (interface{}, error)
	SetCard(ctx context.Context, value models.AnalysedData) error
}

// RedisService caches and retrieves analysed financial data through a provider.
type RedisService struct {
	log      *slog.Logger
	provider AnasyledDataProvider
	TokenTTL time.Duration
}

// NewRedisService creates a new RedisService instance.
func NewRedisService(log *slog.Logger, analyseDataProvider AnasyledDataProvider, tokenTTL time.Duration) *RedisService {
	return &RedisService{
		log:      log,
		TokenTTL: tokenTTL,
		provider: analyseDataProvider,
	}
}

// SetAnalysedData stores analysed data in the cache.
func (r *RedisService) SetAnalysedData(ctx context.Context, dataTitle string, analysedData models.AnalysedData) error {
	const op = "Cache_Service.internal.services.analyseDataProvider.SetAnalysedData"

	err := r.provider.SetCard(ctx, analysedData)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	r.log.Info("SetAnalysedData", slog.String("dataTitle", dataTitle))
	return nil
}

// GetAnalysedData retrieves analysed data from the cache by title.
func (r *RedisService) GetAnalysedData(ctx context.Context, dataTitle string) (interface{}, error) {
	data, err := r.provider.GetValue(ctx, dataTitle)
	if err != nil {
		r.log.Error("Failed to get analysedData", slog.Any("error", err))
		return nil, err
	}
	r.log.Info("GetAnalysedData", slog.String("dataTitle", dataTitle))
	return data, err
}
// Package redis implements Redis-backed storage for the Cache Service,
// handling analysed data storage with sorted-set indexing by channel and timestamp.
package redis

import (
	"CacheService/internal/domain/models"
	"CacheService/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisController manages Redis operations for storing and retrieving analysed data.
type RedisController struct {
	redis      *redis.Client
	DefaultTTl time.Duration
}

// NewRedisController creates a new RedisController with the given connection config.
func NewRedisController(addr string, password string, db int, protocol int, ttl time.Duration) *RedisController {
	return &RedisController{
		redis: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
			Protocol: protocol,
		}),
		DefaultTTl: ttl,
	}
}

// SetCard stores analysed data in Redis as a JSON string and indexes it
// in sorted sets by timestamp and channel username.
func (r *RedisController) SetCard(ctx context.Context, value models.AnalysedData) error {
	const op = "Cache_Service.internal.storage.redis.SetCard"
	id := uuid.NewString()
	timestamp := value.Date.Unix()

	pipeline := r.redis.TxPipeline()

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = pipeline.Set(ctx, "post:"+id, data, r.DefaultTTl).Err()

	if err != nil {
		return err
	}

	pipeline.ZAdd(ctx, "cards", redis.Z{
		Score:  float64(timestamp),
		Member: id,
	})

	pipeline.ZAdd(ctx, "post:channel:"+value.ChannelUsername, redis.Z{
		Score:  float64(timestamp),
		Member: id,
	})

	_, err = pipeline.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// SetValue stores a JSON-serialized value in Redis with optional TTL.
func (r *RedisController) SetValue(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	const op = "Cahce_Service.internal.storage.redis.SetValue"
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(ttl) > 0 {
		return r.redis.Set(ctx, key, data, ttl[0]).Err()
	}
	return r.redis.Set(ctx, key, data, r.DefaultTTl).Err()
}

// GetValue retrieves raw bytes from Redis by key. Returns ErrCacheMiss if not found.
func (r *RedisController) GetValue(ctx context.Context, key string) (interface{}, error) {
	const op = "Cahce_Service.internal.storage.redis.GetValue"
	data, err := r.redis.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, storage.ErrCacheMiss
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return data, nil
}

// Ping checks the Redis connection health.
func (r *RedisController) Ping(ctx context.Context) error {
	return r.redis.Ping(ctx).Err()
}

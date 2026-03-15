package redis

import (
	"CacheService/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisController struct {
	redis      *redis.Client
	DefaultTTl time.Duration
}

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

func (r *RedisController) Ping(ctx context.Context) error {
	return r.redis.Ping(ctx).Err()
}

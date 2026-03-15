package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"recommendationService/internal/storage"
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

func (r *RedisController) GetAll(ctx context.Context, pattern string) (interface{}, error) {
	const op = "Cache_Service.internal.storage.redis.GetAll"

	keys, _, err := r.redis.Scan(ctx, 0, pattern, 10).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	values, err := r.redis.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	log.Println(values)
	return values, nil
}

func (r *RedisController) Ping(ctx context.Context) error {
	return r.redis.Ping(ctx).Err()
}

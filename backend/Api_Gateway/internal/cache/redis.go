package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type RedisCache struct {
	*redis.Client
	DefaultTTl time.Duration
}

func NewRedisCache(addr string, password string, db int, protocol int, ttl time.Duration) *RedisCache {

	return &RedisCache{
		Client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
			Protocol: protocol,
		}),
		DefaultTTl: ttl,
	}
}

func (r *RedisCache) SetValue(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	TTL := r.DefaultTTl
	if len(ttl) > 0 {
		TTL = ttl[0]
	}

	return r.Client.Set(ctx, key, data, TTL).Err()
}

func (r *RedisCache) GetValue(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"recommendationService/internal/domain/models"
	"time"

	"github.com/google/uuid"
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

func (r *RedisController) GetPostsByChannels(ctx context.Context, channels []string) ([]models.Post, error) {
	const op = "Recommendation_Service.internal.storage.redis.GetPostsByChannelsAndDate"

	tmpKey := "post:" + uuid.NewString()
	defer r.redis.Del(ctx, tmpKey)

	postKeys := make([]string, 0, len(channels))
	for _, channel := range channels {
		postKeys = append(postKeys, "post:channel:"+channel)
	}

	_, err := r.redis.ZUnionStore(ctx, tmpKey, &redis.ZStore{Keys: postKeys}).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ids, err := r.redis.ZRevRangeByScore(ctx, tmpKey, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   "+inf",
		Count: 10,
	}).Result()
	if err != nil {
		return []models.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = "post:" + id
	}

	vals, err := r.redis.MGet(ctx, keys...).Result()
	if err != nil {
		return []models.Post{}, fmt.Errorf("%s: %w", op, err)
	}
	posts := make([]models.Post, 0, len(vals))

	for _, val := range vals {
		if val == nil {
			continue
		}
		var post models.Post
		err := json.Unmarshal([]byte(val.(string)), &post)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *RedisController) Ping(ctx context.Context) error {
	return r.redis.Ping(ctx).Err()
}

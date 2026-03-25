package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
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

func (r *RedisController) GetPostsByChannels(ctx context.Context, channels []string, userID int64) ([]models.Post, error) {
	const op = "Recommendation_Service.internal.storage.redis.GetPostsByChannelsAndDate"

	feedKey := fmt.Sprintf("feed:%d", userID)
	cursorKey := fmt.Sprintf("feed_cursor:%d", userID)
	exists, err := r.redis.Exists(ctx, feedKey).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if exists == 0 {
		postKeys := make([]string, 0, len(channels))
		for _, channel := range channels {
			postKeys = append(postKeys, fmt.Sprintf("post:channel:%s", channel))
		}
		_, err := r.redis.ZUnionStore(ctx, feedKey, &redis.ZStore{
			Keys: postKeys,
		}).Result()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		r.redis.Expire(ctx, feedKey, r.DefaultTTl)
	}

	lastScore, err := r.redis.Get(ctx, cursorKey).Float64()
	if err == redis.Nil {
		lastScore = math.Inf(1)
	} else if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	ids, err := r.redis.ZRevRangeByScore(ctx, feedKey, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    fmt.Sprint("%f", lastScore),
		Offset: 0,
		Count:  2,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(ids) == 0 {
		return []models.Post{}, nil
	}
	lastId := ids[len(ids)-1]
	score, err := r.redis.ZScore(ctx, feedKey, lastId).Result()
	if err == nil {
		r.redis.Set(ctx, cursorKey, score, time.Hour)
	}

	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, fmt.Sprintf("post:%s", id))
	}
	vals, err := r.redis.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	posts := make([]models.Post, 0, len(keys))
	for _, val := range vals {
		if val == nil {
			continue
		}
		var post models.Post
		err := json.Unmarshal([]byte(val.(string)), &post)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *RedisController) Ping(ctx context.Context) error {
	return r.redis.Ping(ctx).Err()
}

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"recommendationService/internal/domain/models"
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

func (r *RedisController) GetPostsByChannels(ctx context.Context, channels []string, userID int64, cursor *models.Cursor, limit int64) ([]models.Post, *models.Cursor, error) {
	const op = "Recommendation_Service.internal.storage.redis.GetPostsByChannelsAndDate"

	tmpKey := fmt.Sprintf(
		"tmp:feed:%d:%d",
		userID,
		time.Now().UnixNano(),
	)

	if err := r.initilizeNewZUnion(ctx, tmpKey, channels); err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	maxScore := "+inf"
	if cursor != nil {
		maxScore = fmt.Sprintf("(%f", cursor.Score)
	}

	zPosts, err := r.getSortedPosts(ctx, tmpKey, maxScore, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(zPosts) == 0 {
		return []models.Post{}, nil, nil
	}

	postKeys := preparePostKeys(zPosts, cursor)
	if len(postKeys) == 0 {
		return []models.Post{}, nil, nil
	}

	vals, err := r.redis.MGet(ctx, postKeys...).Result()
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	posts := unmarshalToPosts(vals)
	nextCursor, err := configureNewCursor(zPosts)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nextCursor, nil
}

func unmarshalToPosts(vals []any) []models.Post {
	posts := make([]models.Post, 0, len(vals))
	for _, val := range vals {
		if val == nil {
			continue
		}
		strId, ok := val.(string)
		if !ok {
			continue
		}
		var post models.Post
		if err := json.Unmarshal([]byte(strId), &post); err != nil {
			continue
		}
		posts = append(posts, post)

	}
	return posts
}

func (r *RedisController) initilizeNewZUnion(ctx context.Context, tmpKey string, channels []string) error {
	const op = "Recommendation_Service.internal.storage.redis.inicilizeNewZUnion"
	channelKeys := make([]string, 0, len(channels))
	for _, ch := range channels {
		channelKeys = append(channelKeys, fmt.Sprintf("post:channel:%s", ch))
	}
	if err := r.redis.ZUnionStore(
		ctx,
		tmpKey,
		&redis.ZStore{Keys: channelKeys},
	).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	r.redis.Expire(ctx, tmpKey, 5*time.Second)
	return nil
}

func configureNewCursor(zPosts []redis.Z) (*models.Cursor, error) {
	const op = "Recommendation_Service.internal.storage.redis.configureNewCursor"

	last := zPosts[len(zPosts)-1]
	lastID, ok := last.Member.(string)
	if !ok {
		return nil, fmt.Errorf("%s: can not get new Cursor", op)
	}
	nextCursor := &models.Cursor{
		Score: last.Score,
		ID:    lastID,
	}
	return nextCursor, nil
}

func preparePostKeys(zPosts []redis.Z, cursor *models.Cursor) []string {
	postKeys := make([]string, 0, len(zPosts))
	for _, post := range zPosts {
		strId, ok := post.Member.(string)
		if !ok {
			continue
		}
		if cursor != nil && post.Score == cursor.Score && strId == cursor.ID {
			continue
		}
		postKeys = append(postKeys, "post:"+strId)
	}
	return postKeys
}

func (r *RedisController) getSortedPosts(ctx context.Context, tmpKey string, maxScore string, count int64) ([]redis.Z, error) {
	posts, err := r.redis.ZRevRangeByScoreWithScores(
		ctx,
		tmpKey,
		&redis.ZRangeBy{
			Min:   "-inf",
			Max:   maxScore,
			Count: count,
		},
	).Result()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *RedisController) Ping(ctx context.Context) error {
	return r.redis.Ping(ctx).Err()
}

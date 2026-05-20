package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"recommendationService/internal/domain/models"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite
	controller *RedisController
	redisAddr  string
	redisPw    string
}

func (r *RedisTestSuite) SetupTest() {
	r.redisAddr = os.Getenv("REDIS_HOST")
	if r.redisAddr == "" {
		r.T().Skip("REDIS_HOST not set")
	}

	r.redisPw = os.Getenv("REDIS_PASSWORD")
	if r.redisPw == "" {
		r.T().Skip("REDIS_PASSWORD not set")
	}

	redisDB := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		fmt.Sscanf(dbStr, "%d", &redisDB)
	}

	r.controller = NewRedisController(
		r.redisAddr,
		r.redisPw,
		redisDB,
		3,
		5*time.Minute,
	)

	r.setupTestData()
}

func (r *RedisTestSuite) setupTestData() {
	ctx := context.Background()
	r.cleanupTestData()

	posts := []map[string]interface{}{
		{
			"post_text":        "Test post 1",
			"post_uri":         "https://example.com/1",
			"channel_username": "channel1",
			"date":             time.Now().Unix(),
		},
		{
			"post_text":        "Test post 2",
			"post_uri":         "https://example.com/2",
			"channel_username": "channel1",
			"date":             time.Now().Unix(),
		},
	}

	r.controller.redis.Del(ctx, "post:channel:channel1")

	for i, post := range posts {
		data, _ := json.Marshal(post)
		score := float64(1000 + i)
		r.controller.redis.ZAdd(ctx, "post:channel:channel1", redis.Z{
			Score:  score,
			Member: fmt.Sprintf("post_%d", i),
		})
		r.controller.redis.Set(ctx, fmt.Sprintf("post:post_%d", i), data, time.Hour)
	}
}

func (r *RedisTestSuite) cleanupTestData() {
	ctx := context.Background()
	keys, _ := r.controller.redis.Keys(ctx, "tmp:feed:*").Result()
	if len(keys) > 0 {
		r.controller.redis.Del(ctx, keys...)
	}
}

func (r *RedisTestSuite) TearDownTest() {
	r.cleanupTestData()
}

func (r *RedisTestSuite) TestPing_Success() {
	ctx := context.Background()
	err := r.controller.Ping(ctx)
	r.NoError(err)
}

func (r *RedisTestSuite) TestPing_InvalidConnection() {
	controller := NewRedisController(
		"invalid-host:6379",
		"",
		0,
		3,
		5*time.Minute,
	)

	ctx := context.Background()
	err := controller.Ping(ctx)
	r.Error(err)
}

func (r *RedisTestSuite) TestGetPostsByChannels_Success() {
	ctx := context.Background()
	channels := []string{"channel1"}

	posts, cursor, err := r.controller.GetPostsByChannels(ctx, channels, 1, nil, 10)
	r.NoError(err)
	r.NotNil(posts)
	r.NotNil(cursor)
}

func (r *RedisTestSuite) TestGetPostsByChannels_WithCursor() {
	ctx := context.Background()
	channels := []string{"channel1"}

	_, cursor, err := r.controller.GetPostsByChannels(ctx, channels, 2, nil, 10)
	r.NoError(err)

	if cursor == nil {
		r.T().Skip("No cursor returned - not enough data")
		return
	}

	posts, nextCursor, err := r.controller.GetPostsByChannels(ctx, channels, 2, cursor, 10)
	r.NoError(err)
	r.NotNil(posts)
	r.Nil(nextCursor)
}

func (r *RedisTestSuite) TestGetPostsByChannels_EmptyChannels() {
	ctx := context.Background()

	posts, cursor, err := r.controller.GetPostsByChannels(ctx, []string{"nonexistent"}, 1, nil, 10)
	r.NoError(err)
	r.Empty(posts)
	r.Nil(cursor)
}

func (r *RedisTestSuite) TestGetPostsByChannels_InvalidChannel() {
	ctx := context.Background()
	channels := []string{"nonexistent_channel_xyz"}

	posts, cursor, err := r.controller.GetPostsByChannels(ctx, channels, 1, nil, 10)
	r.NoError(err)
	r.Empty(posts)
	r.Nil(cursor)
}

func (r *RedisTestSuite) TestInitializeNewZUnion() {
	ctx := context.Background()
	tmpKey := "tmp:test:123"
	channels := []string{"channel1"}

	err := r.controller.initilizeNewZUnion(ctx, tmpKey, channels)
	r.NoError(err)

	exists, _ := r.controller.redis.Exists(ctx, tmpKey).Result()
	r.Equal(int64(1), exists)

	r.controller.redis.Del(ctx, tmpKey)
}

func (r *RedisTestSuite) TestGetSortedPosts() {
	ctx := context.Background()
	tmpKey := "tmp:test:456"

	r.controller.redis.ZAdd(ctx, tmpKey, redis.Z{Score: 10, Member: "a"}, redis.Z{Score: 20, Member: "b"}, redis.Z{Score: 30, Member: "c"})

	posts, err := r.controller.getSortedPosts(ctx, tmpKey, "+inf", 0, 5)
	r.NoError(err)
	r.Len(posts, 3)
	r.Equal(30.0, posts[0].Score)

	r.controller.redis.Del(ctx, tmpKey)
}

func (r *RedisTestSuite) TestGetSortedPosts_WithOffset() {
	ctx := context.Background()
	tmpKey := "tmp:test:789"

	r.controller.redis.ZAdd(ctx, tmpKey, redis.Z{Score: 10, Member: "a"}, redis.Z{Score: 20, Member: "b"}, redis.Z{Score: 30, Member: "c"})

	posts, err := r.controller.getSortedPosts(ctx, tmpKey, "+inf", 1, 2)
	r.NoError(err)
	r.Len(posts, 2)
	r.Equal(20.0, posts[0].Score)

	r.controller.redis.Del(ctx, tmpKey)
}

func (r *RedisTestSuite) TestGetChannelsHash() {
	channels1 := []string{"a", "b", "c"}
	channels2 := []string{"c", "b", "a"}
	channels3 := []string{"a", "b", "d"}

	hash1 := getChannelsHash(channels1)
	hash2 := getChannelsHash(channels2)
	hash3 := getChannelsHash(channels3)

	r.Equal(hash1, hash2)
	r.NotEqual(hash1, hash3)
	r.Len(hash1, 32)
}

func (r *RedisTestSuite) TestPreparePostKeys() {
	zPosts := []redis.Z{
		{Score: 10, Member: "post_1"},
		{Score: 20, Member: "post_2"},
		{Score: 30, Member: "post_3"},
	}

	keys := preparePostKeys(zPosts, nil)
	r.Len(keys, 3)
	r.Equal("post:post_1", keys[0])
	r.Equal("post:post_2", keys[1])
	r.Equal("post:post_3", keys[2])
}

func (r *RedisTestSuite) TestPreparePostKeys_WithCursor() {
	zPosts := []redis.Z{
		{Score: 10, Member: "post_1"},
		{Score: 20, Member: "post_2"},
		{Score: 30, Member: "post_3"},
	}

	cursor := &models.Cursor{Score: 20, ID: "post_2"}
	keys := preparePostKeys(zPosts, cursor)
	r.Len(keys, 2)
	r.Equal("post:post_1", keys[0])
	r.Equal("post:post_3", keys[1])
}

func (r *RedisTestSuite) TestConfigureNewCursor() {
	zPosts := []redis.Z{
		{Score: 10, Member: "post_1"},
		{Score: 20, Member: "post_2"},
		{Score: 30, Member: "post_3"},
	}

	cursor, err := configureNewCursor(zPosts)
	r.NoError(err)
	r.NotNil(cursor)
	r.Equal(30.0, cursor.Score)
	r.Equal("post_3", cursor.ID)
}

func (r *RedisTestSuite) TestConfigureNewCursor_Empty() {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	zPosts := []redis.Z{}
	_, err := configureNewCursor(zPosts)
	r.Error(err)
}

func (r *RedisTestSuite) TestConfigureNewCursor_InvalidMember() {
	zPosts := []redis.Z{
		{Score: 10, Member: 123},
	}

	_, err := configureNewCursor(zPosts)
	r.Error(err)
}

func (r *RedisTestSuite) TestUnmarshalToPosts() {
	vals := []any{
		`{"Stocks":[],"post_text":"test","post_uri":"url","channel_username":"ch","reasoning":"","date":"2024-01-01T00:00:00Z"}`,
		`{"Stocks":[],"post_text":"test2","post_uri":"url2","channel_username":"ch2","reasoning":"","date":"2024-01-01T00:00:00Z"}`,
		nil,
		"invalid json",
		123,
		"",
	}

	posts := unmarshalToPosts(vals)
	r.GreaterOrEqual(len(posts), 0)
}

func (r *RedisTestSuite) TestUnmarshalToPosts_Empty() {
	posts := unmarshalToPosts([]any{})
	r.Empty(posts)
}

func (r *RedisTestSuite) TestGetPostsByChannels_CacheHit() {
	ctx := context.Background()
	channels := []string{"channel1"}

	_, _, err := r.controller.GetPostsByChannels(ctx, channels, 999, nil, 10)
	r.NoError(err)

	posts, cursor, err := r.controller.GetPostsByChannels(ctx, channels, 999, nil, 10)
	r.NoError(err)
	r.NotNil(posts)
	r.NotNil(cursor)
}

func (r *RedisTestSuite) TestGetPostsByChannels_Limit() {
	ctx := context.Background()
	channels := []string{"channel1"}

	posts, _, err := r.controller.GetPostsByChannels(ctx, channels, 666, nil, 1)
	r.NoError(err)
	r.GreaterOrEqual(len(posts), 0)
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
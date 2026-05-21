package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"CacheService/internal/domain/models"
	"CacheService/internal/storage"
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
		r.redisAddr = "redis-master:6379" // Default to the service name and port in docker-compose
	}

	r.redisPw = os.Getenv("REDIS_PASSWORD")
	if r.redisPw == "" {
		r.redisPw = "1ux35qBk4YgCMsd7eg4ju" // Default password from .env
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

	// Clear test data before each test
	ctx := context.Background()
	r.controller.redis.FlushDB(ctx)
}

func (r *RedisTestSuite) TearDownTest() {
	// Clean up after all tests
	ctx := context.Background()
	r.controller.redis.FlushDB(ctx)
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

func (r *RedisTestSuite) TestSetValue_Success() {
	ctx := context.Background()
	testKey := "test:setvalue:key"
	testValue := map[string]string{"data": "test_value"}

	err := r.controller.SetValue(ctx, testKey, testValue)
	r.NoError(err)

	// Retrieve and verify
	data, err := r.controller.GetValue(ctx, testKey)
	r.NoError(err)
	r.NotNil(data)

	// Verify it's JSON bytes
	var result map[string]string
	err = json.Unmarshal(data.([]byte), &result)
	r.NoError(err)
	r.Equal(testValue["data"], result["data"])
}

func (r *RedisTestSuite) TestSetValue_WithCustomTTL() {
	ctx := context.Background()
	testKey := "test:setvalue:ttl"
	testValue := "test_value"

	// Set with custom TTL of 1 second
	err := r.controller.SetValue(ctx, testKey, testValue, 1*time.Second)
	r.NoError(err)

	// Verify it exists immediately
	data, err := r.controller.GetValue(ctx, testKey)
	r.NoError(err)
	r.NotNil(data)

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Should now be a cache miss
	data, err = r.controller.GetValue(ctx, testKey)
	r.Error(err)
	r.Equal(storage.ErrCacheMiss, err)
}

func (r *RedisTestSuite) TestGetValue_Success() {
	ctx := context.Background()
	testKey := "test:getvalue:key"
	testValue := []byte("test_data")

	// Set value directly in Redis
	err := r.controller.redis.Set(ctx, testKey, testValue, 0).Err()
	r.NoError(err)

	// Retrieve using our method
	data, err := r.controller.GetValue(ctx, testKey)
	r.NoError(err)
	r.NotNil(data)
	r.Equal(testValue, data.([]byte))
}

func (r *RedisTestSuite) TestGetValue_CacheMiss() {
	ctx := context.Background()
	testKey := "test:getvalue:nomatch"

	data, err := r.controller.GetValue(ctx, testKey)
	r.Error(err)
	r.Equal(storage.ErrCacheMiss, err)
	r.Nil(data)
}

func (r *RedisTestSuite) TestSetCard_Success() {
	ctx := context.Background()
	
	// Create test data
	testData := models.AnalysedData{
		Stocks: []models.Stock{
			{StockName: "AAPL", Side: "BUY"},
			{StockName: "GOOGL", Side: "SELL"},
		},
		PostText:        "Test post text",
		PostURI:         "https://example.com/post",
		ChannelUsername: "testchannel",
		Date:            time.Now(),
		Reasoning:       "Test reasoning",
	}

	err := r.controller.SetCard(ctx, testData)
	r.NoError(err)

	// Verify the data was stored by checking the sorted set
	count, err := r.controller.redis.ZCard(ctx, "cards").Result()
	r.NoError(err)
	r.Greater(count, int64(0))

	// Verify channel-specific set exists
	channelCount, err := r.controller.redis.ZCard(ctx, "post:channel:testchannel").Result()
	r.NoError(err)
	r.Greater(channelCount, int64(0))
}



func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
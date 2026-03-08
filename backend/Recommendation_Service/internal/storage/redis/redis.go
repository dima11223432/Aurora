package redis

import "github.com/redis/go-redis/v9"

type RedisController struct {
	redis *redis.Client
}

package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{rdb: rdb}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.rdb.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

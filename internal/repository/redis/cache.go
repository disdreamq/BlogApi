package redis

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb       *redis.Client
	cacheMiss atomic.Uint64
	cacheHit  atomic.Uint64
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{rdb: rdb}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.rdb.Set(ctx, key, value.([]byte), ttl).Err()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

func (r *RedisCache) ShowRatio() (float64, error) {
	miss := r.cacheMiss.Load()
	hit := r.cacheHit.Load()
	if miss == 0 {
		return float64(hit), errors.New("Cache miss is 0")
	}
	return float64(hit) / float64(miss), nil

}

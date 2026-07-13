package redis

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisCache struct {
	rdb       *redis.Client
	cacheMiss atomic.Uint64
	cacheHit  atomic.Uint64
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{rdb: rdb}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, bool) {
	val, err := r.rdb.Get(ctx, key).Result()
	logger := log.Ctx(ctx)
	if err != nil {
		switch err {
		case redis.Nil:
			return val, false
		default:
			logger.Err(err).
				Str("trace_id", ctx.Value("trace_id").(string)).
				Str("key", key)
		}
	}

	hit := val != ""
	if hit {
		r.cacheMiss.Add(1)
	} else {
		r.cacheHit.Add(1)
	}
	logger.Debug().
		Str("trace_id", ctx.Value("trace_id").(string)).
		Str("key", key).
		Bool("hit", hit)
	return val, true
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) bool {
	logger := log.Ctx(ctx)
	err := r.rdb.Set(ctx, key, value.([]byte), ttl).Err()
	if err != nil {
		logger.Err(err).
			Str("key", key)
		return false
	}
	logger.Debug().
		Str("trace_id", ctx.Value("trace_id").(string)).
		Str("key", key).
		Msg("Set to cache")
	return true
}

func (r *RedisCache) Del(ctx context.Context, key string) bool {
	logger := log.Ctx(ctx)
	err := r.rdb.Del(ctx, key).Err()
	if err != nil {
		logger.Err(err).
			Str("trace_id", ctx.Value("trace_id").(string)).
			Str("key", key)
		return false
	}
	logger.Debug().
		Str("trace_id", ctx.Value("trace_id").(string)).
		Str("key", key).
		Msg("Deleted from cache")
	return true
}

func (r *RedisCache) ShowRatio() (float64, error) {
	miss := r.cacheMiss.Load()
	hit := r.cacheHit.Load()
	if miss == 0 {
		return float64(hit), errors.New("Cache miss is 0")
	}
	return float64(hit) / float64(miss), nil

}

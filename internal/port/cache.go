package port

import (
	"context"
	"time"
)

type CacheGetter interface {
	Get(ctx context.Context, key string) (string, error)
}

type CacheSetter interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
}

type CacheDeleter interface {
	Del(ctx context.Context, key string) error
}

type Cache interface {
	CacheGetter
	CacheSetter
	CacheDeleter
}

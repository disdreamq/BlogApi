package port

import (
	"context"
	"time"
)

type CacheGetter interface {
	Get(ctx context.Context, key string) (string, bool)
}

type CacheSetter interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) bool
}

type CacheDeleter interface {
	Del(ctx context.Context, key string) bool
}

type Cache interface {
	CacheGetter
	CacheSetter
	CacheDeleter
}

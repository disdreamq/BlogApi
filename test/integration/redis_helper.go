package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	r "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestCache struct {
	Container *r.RedisContainer
	Addr      string
	Client    *redis.Client
}

func StartRedisContainer(ctx context.Context) (*TestCache, error) {
	redisContainer, err := r.Run(ctx,
		"redis:8.8-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp"),
			wait.ForLog("Ready to accept connections").WithStartupTimeout(180*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start redis container: %w", err)
	}

	host, err := redisContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %w", err)
	}
	port, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get port: %w", err)
	}

	if host == "localhost" {
		host = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%s", host, port.Port())

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   6, // 7
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &TestCache{
		Container: redisContainer,
		Addr:      addr,
		Client:    client,
	}, nil
}

func StartTestRedis(t *testing.T) *TestCache {
	testRedis, err := StartRedisContainer(context.Background())
	if err != nil {
		t.Fatal("failed to start redis container:", err)
	}

	testcontainers.CleanupContainer(t, testRedis.Container)

	return testRedis
}

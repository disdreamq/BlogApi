package integration

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	Container *postgres.PostgresContainer
	ConnStr   string
	DB        *sqlx.DB
}

func StartTestPostgres(t *testing.T) *TestDatabase {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatal("failed to start postgres container:", err)
	}

	testcontainers.CleanupContainer(t, pgContainer)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal("failed to get connection string:", err)
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Fatal("failed to open database connection:", err)
	}

	if err := db.PingContext(ctx); err != nil {
		t.Fatal("failed to ping database:", err)
	}

	return &TestDatabase{
		Container: pgContainer,
		ConnStr:   connStr,
		DB:        db,
	}
}

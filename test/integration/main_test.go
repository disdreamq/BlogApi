package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *TestDatabase

func TestMain(m *testing.M) {
	testDB = startTestDatabase()

	code := m.Run()

	if testDB.DB != nil {
		testDB.DB.Close()
	}

	os.Exit(code)
}

func startTestDatabase() *TestDatabase {
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
		panic("failed to start postgres container: " + err.Error())
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic("failed to get connection string: " + err.Error())
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic("failed to open database connection: " + err.Error())
	}

	return &TestDatabase{
		Container: pgContainer,
		ConnStr:   connStr,
		DB:        db,
	}
}

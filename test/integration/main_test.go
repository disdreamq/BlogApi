package integration

import (
	"context"
	"os"
	"testing"

	"github.com/disdreamq/BlogApi/test/fixtures"
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
	testDB, err := StartPostgresContainer(ctx)
	if err != nil {
		panic("failed to start postgres container: " + err.Error())
	}

	if err := fixtures.RunMigrations(testDB.ConnStr); err != nil {
		panic(err)
	}

	return testDB
}

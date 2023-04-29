package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	postgres *Postgres
)

func testDatabaseURL() (string, bool) {
	return os.LookupEnv("TEST_DATABASE_URL")
}

func shouldSkip() bool {
	_, exists := testDatabaseURL()
	return !exists
}

func setup(t *testing.T) {
	t.Helper()

	gofakeit.Seed(0)

	if postgres == nil {
		pg, err := newTestPostgres(t)
		if err != nil {
			t.Fatal(err)
		}

		postgres = pg
	}
}

func newTestPostgres(t *testing.T) (*Postgres, error) {
	databaseURL, exists := testDatabaseURL()
	if !exists {
		t.Fatal("failed to lookup database URL")
	}

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	pg, err := Connect(context.Background(), cfg, time.Now)
	if err != nil {
		return nil, err
	}

	return pg, nil
}

package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	if shouldSkip() {
		t.SkipNow()
	}

	pg, err := newTestPostgres(t)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, pg)
	assert.NotNil(t, pg.db)
}

func TestClose(t *testing.T) {
	if shouldSkip() {
		t.SkipNow()
	}

	pg, err := newTestPostgres(t)
	if err != nil {
		t.Error(err)
	}

	assert.NotPanics(t, func() { pg.Close() })
}

func TestCheck(t *testing.T) {
	if shouldSkip() {
		t.SkipNow()
	}

	pg, err := newTestPostgres(t)
	if err != nil {
		t.Error(err)
	}

	if err := pg.Check(context.Background()); err != nil {
		t.Error(err)
	}
}

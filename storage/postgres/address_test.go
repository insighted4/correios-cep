package postgres

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/insighted4/correios-cep/pkg/errors"
	"github.com/insighted4/correios-cep/storage"
	"github.com/stretchr/testify/assert"
)

func TestPostgres_Create(t *testing.T) {
	if shouldSkip() {
		t.SkipNow()
	}

	setup(t)

	p1 := &storage.Address{
		CEP:          gofakeit.UUID(),
		State:        gofakeit.LoremIpsumWord(),
		City:         gofakeit.LoremIpsumWord(),
		Neighborhood: gofakeit.LoremIpsumWord(),
		Location:     gofakeit.LoremIpsumWord(),
		Source:       gofakeit.Categories(),
	}

	ctx := context.Background()
	if err := postgres.Create(ctx, p1); err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, p1.CreatedAt)
	assert.NotNil(t, p1.UpdatedAt)

	p2, err := postgres.get(ctx, p1.CEP, errors.Op("TestPostgres_Create"))
	if err != nil {
		t.Fatal(err)
	}

	j1, err := json.Marshal(p1)
	if err != nil {
		t.Fatal(err)
	}

	j2, err := json.Marshal(p2)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(j1), string(j2))
}

func TestPostgres_Update(t *testing.T) {
	if shouldSkip() {
		t.SkipNow()
	}

	setup(t)

	p0 := &storage.Address{
		CEP:          gofakeit.UUID(),
		State:        gofakeit.LoremIpsumWord(),
		City:         gofakeit.LoremIpsumWord(),
		Neighborhood: gofakeit.LoremIpsumWord(),
		Location:     gofakeit.LoremIpsumWord(),
		Source:       gofakeit.Categories(),
	}

	ctx := context.Background()
	if err := postgres.Create(ctx, p0); err != nil {
		t.Fatal(err)
	}

	p1, err := postgres.get(ctx, p0.CEP, errors.Op("TestPostgres_UpdateProduct"))
	if err != nil {
		t.Fatal(err)
	}

	p2UID := gofakeit.UUID()
	updater := func(old *storage.Address) (*storage.Address, error) {
		old.CEP = p2UID
		old.State = gofakeit.LoremIpsumWord()
		old.City = gofakeit.LoremIpsumWord()
		old.Neighborhood = gofakeit.LoremIpsumWord()
		old.Location = gofakeit.LoremIpsumWord()
		old.Source = gofakeit.Categories()
		return old, nil
	}

	if err := postgres.Update(ctx, p0.CEP, updater); err != nil {
		t.Fatal(err)
	}

	p2, err := postgres.get(ctx, p2UID, errors.Op("TestPostgres_UpdateProduct"))
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, p1.CEP, p2.CEP)
	assert.NotEqual(t, p1.State, p2.State)
	assert.NotEqual(t, p1.City, p2.City)
	assert.NotEqual(t, p1.Neighborhood, p2.Neighborhood)
	assert.NotEqual(t, p1.Location, p2.Location)
	assert.NotEqual(t, p1.Source, p2.Source)
	assert.Equal(t, p1.CreatedAt.UTC(), p2.CreatedAt.UTC())
	assert.GreaterOrEqual(t, p2.UpdatedAt.UTC(), p1.UpdatedAt.UTC())
}

func TestPostgres_GetNotFound(t *testing.T) {
	if shouldSkip() {
		t.SkipNow()
	}

	setup(t)

	p, err := postgres.Get(context.Background(), gofakeit.UUID())
	assert.True(t, errors.Is(err, errors.KindNotFound))
	assert.Nil(t, p)
}

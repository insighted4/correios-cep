// Copyright 2023 The Correios CEP Admin Authors
//
// Licensed under the AGPL, Version 3.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.gnu.org/licenses/agpl-3.0.en.html
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgres

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/insighted4/correios-cep/pkg/errors"
	"github.com/insighted4/correios-cep/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		Children: []*storage.Address{
			{
				CEP:          gofakeit.UUID(),
				State:        gofakeit.LoremIpsumWord(),
				City:         gofakeit.LoremIpsumWord(),
				Neighborhood: gofakeit.LoremIpsumWord(),
				Location:     gofakeit.LoremIpsumWord(),
			},
		},
	}

	ctx := context.Background()
	err := postgres.CreateAddress(ctx, p1)
	require.NoError(t, err)

	assert.NotNil(t, p1.CreatedAt)
	assert.NotNil(t, p1.UpdatedAt)

	p2, err := postgres.get(ctx, p1.CEP, errors.Op("TestPostgres_Create"))
	require.NoError(t, err)

	j1, err := json.Marshal(p1)
	require.NoError(t, err)

	j2, err := json.Marshal(p2)
	require.NoError(t, err)

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
		Children: []*storage.Address{
			{
				CEP:          gofakeit.UUID(),
				State:        gofakeit.LoremIpsumWord(),
				City:         gofakeit.LoremIpsumWord(),
				Neighborhood: gofakeit.LoremIpsumWord(),
				Location:     gofakeit.LoremIpsumWord(),
			},
		},
	}

	ctx := context.Background()
	err := postgres.CreateAddress(ctx, p0)
	require.NoError(t, err)

	p1, err := postgres.get(ctx, p0.CEP, errors.Op("TestPostgres_UpdateProduct"))
	require.NoError(t, err)

	p2UID := gofakeit.UUID()
	updater := func(old *storage.Address) (*storage.Address, error) {
		old.CEP = p2UID
		old.State = gofakeit.LoremIpsumWord()
		old.City = gofakeit.LoremIpsumWord()
		old.Neighborhood = gofakeit.LoremIpsumWord()
		old.Location = gofakeit.LoremIpsumWord()
		old.Children = []*storage.Address{
			{
				CEP:          gofakeit.UUID(),
				State:        gofakeit.LoremIpsumWord(),
				City:         gofakeit.LoremIpsumWord(),
				Neighborhood: gofakeit.LoremIpsumWord(),
				Location:     gofakeit.LoremIpsumWord(),
			},
		}
		return old, nil
	}

	err = postgres.UpdateAddress(ctx, p0.CEP, updater)
	require.NoError(t, err)

	p2, err := postgres.get(ctx, p2UID, errors.Op("TestPostgres_UpdateProduct"))
	require.NoError(t, err)

	assert.NotEqual(t, p1.CEP, p2.CEP)
	assert.NotEqual(t, p1.State, p2.State)
	assert.NotEqual(t, p1.City, p2.City)
	assert.NotEqual(t, p1.Neighborhood, p2.Neighborhood)
	assert.NotEqual(t, p1.Location, p2.Location)
	assert.NotEqual(t, p1.Children, p2.Children)
	assert.Equal(t, p1.CreatedAt.UTC(), p2.CreatedAt.UTC())
	assert.GreaterOrEqual(t, p2.UpdatedAt.UTC(), p1.UpdatedAt.UTC())
}

func TestPostgres_GetNotFound(t *testing.T) {
	if shouldSkip() {
		t.SkipNow()
	}

	setup(t)

	p, err := postgres.GetAddress(context.Background(), gofakeit.UUID())
	assert.True(t, errors.Is(err, errors.KindNotFound))
	assert.Nil(t, p)
}

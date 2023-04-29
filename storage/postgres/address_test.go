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
		CEP:                      gofakeit.UUID(),
		UF:                       gofakeit.State(),
		LocalidadeNumero:         gofakeit.StreetNumber(),
		Localidade:               gofakeit.Street(),
		LogradouroDNEC:           gofakeit.UUID(),
		BairoNumero:              gofakeit.UUID(),
		Bairro:                   gofakeit.Word(),
		FaixasCaixaPostal:        []interface{}{gofakeit.NiceColors()},
		FaixasCEP:                []interface{}{gofakeit.NiceColors()},
		LocNoSem:                 gofakeit.UUID(),
		LocalidadeSubordinada:    gofakeit.UUID(),
		LogradouroTexto:          gofakeit.LoremIpsumWord(),
		LogradouroTextoAdicional: gofakeit.LoremIpsumWord(),
		NomeUnidade:              gofakeit.UUID(),
		NumeroLocalidade:         gofakeit.UUID(),
		Situacao:                 gofakeit.UUID(),
		TipoCEP:                  gofakeit.StreetNumber(),
		Desmembramento:           []interface{}{gofakeit.NiceColors()},
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
		CEP:                      gofakeit.UUID(),
		UF:                       gofakeit.State(),
		LocalidadeNumero:         gofakeit.StreetNumber(),
		Localidade:               gofakeit.Street(),
		LogradouroDNEC:           gofakeit.UUID(),
		BairoNumero:              gofakeit.UUID(),
		Bairro:                   gofakeit.Word(),
		FaixasCaixaPostal:        []interface{}{gofakeit.NiceColors()},
		FaixasCEP:                []interface{}{gofakeit.NiceColors()},
		LocNoSem:                 gofakeit.UUID(),
		LocalidadeSubordinada:    gofakeit.UUID(),
		LogradouroTexto:          gofakeit.LoremIpsumWord(),
		LogradouroTextoAdicional: gofakeit.LoremIpsumWord(),
		NomeUnidade:              gofakeit.UUID(),
		NumeroLocalidade:         gofakeit.UUID(),
		Situacao:                 gofakeit.UUID(),
		TipoCEP:                  gofakeit.StreetNumber(),
		Desmembramento:           []interface{}{gofakeit.NiceColors()},
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
		old.UF = gofakeit.State()
		old.LocalidadeNumero = gofakeit.StreetNumber()
		old.Localidade = gofakeit.Street()
		old.LogradouroDNEC = gofakeit.UUID()
		old.BairoNumero = gofakeit.UUID()
		old.Bairro = gofakeit.Word()
		old.FaixasCaixaPostal = []interface{}{gofakeit.NiceColors()}
		old.FaixasCEP = []interface{}{gofakeit.NiceColors()}
		old.LocNoSem = gofakeit.UUID()
		old.LocalidadeSubordinada = gofakeit.UUID()
		old.LogradouroTexto = gofakeit.LoremIpsumWord()
		old.LogradouroTextoAdicional = gofakeit.LoremIpsumWord()
		old.NomeUnidade = gofakeit.UUID()
		old.NumeroLocalidade = gofakeit.UUID()
		old.Situacao = gofakeit.UUID()
		old.TipoCEP = gofakeit.StreetNumber()
		old.Desmembramento = []interface{}{gofakeit.NiceColors()}
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
	assert.NotEqual(t, p1.UF, p2.UF)
	assert.NotEqual(t, p1.LocalidadeNumero, p2.LocalidadeNumero)
	assert.NotEqual(t, p1.Localidade, p2.Localidade)
	assert.NotEqual(t, p1.LogradouroDNEC, p2.LogradouroDNEC)
	assert.NotEqual(t, p1.BairoNumero, p2.BairoNumero)
	assert.NotEqual(t, p1.Bairro, p2.Bairro)
	assert.NotEqual(t, p1.FaixasCaixaPostal, p2.FaixasCaixaPostal)
	assert.NotEqual(t, p1.FaixasCEP, p2.FaixasCEP)
	assert.NotEqual(t, p1.LocNoSem, p2.LocNoSem)
	assert.NotEqual(t, p1.LocalidadeSubordinada, p2.LocalidadeSubordinada)
	assert.NotEqual(t, p1.LogradouroTexto, p2.LogradouroTexto)
	assert.NotEqual(t, p1.LogradouroTextoAdicional, p2.LogradouroTextoAdicional)
	assert.NotEqual(t, p1.NomeUnidade, p2.NomeUnidade)
	assert.NotEqual(t, p1.NumeroLocalidade, p2.NumeroLocalidade)
	assert.NotEqual(t, p1.Situacao, p2.Situacao)
	assert.NotEqual(t, p1.TipoCEP, p2.TipoCEP)
	assert.NotEqual(t, p1.Desmembramento, p2.Desmembramento)
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

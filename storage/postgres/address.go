package postgres

import (
	"context"

	"github.com/insighted4/correios-cep/pkg/errors"
	"github.com/insighted4/correios-cep/storage"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) Create(ctx context.Context, address *storage.Address) error {
	const op errors.Op = "postgres.CreateAddressCode"
	query := `INSERT INTO addresses (
				cep,
				uf,
				localidade_numero,
				localidade,
				logradouro_dnec,
				bairro_numero,
				bairro,
				faixas_caixa_postal,
				faixas_cep,
				loc_no_sem,
				localidade_subordinada,
				logradouro_texto,
				logradouro_texto_adicional,
				nome_unidade,
				numero_localidade,
				situacao,
				tipo_cep,
				desmembramento,
				created_at,
				updated_at
    		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20); 
	`

	now := p.now().UTC()
	address.CreatedAt = &now
	address.UpdatedAt = &now

	if _, err := p.db.Exec(ctx, query,
		address.CEP,
		address.UF,
		address.LocalidadeNumero,
		address.Localidade,
		address.LogradouroDNEC,
		address.BairoNumero,
		address.Bairro,
		address.FaixasCaixaPostal,
		address.FaixasCEP,
		address.LocNoSem,
		address.LocalidadeSubordinada,
		address.LogradouroTexto,
		address.LogradouroTextoAdicional,
		address.NomeUnidade,
		address.NumeroLocalidade,
		address.Situacao,
		address.TipoCEP,
		address.Desmembramento,
		address.CreatedAt,
		address.UpdatedAt,
	); err != nil {
		return errors.E(op, kind(err), err)
	}

	return nil
}

func (p *Postgres) Update(ctx context.Context, cep string, updater storage.Updater) error {
	const op errors.Op = "postgres.UpdateAddressCode"

	updateFn := func(tx pgx.Tx) error {
		old, err := p.get(ctx, cep, op)
		if err != nil {
			return err
		}

		address, err := updater(old)
		if err != nil {
			return err
		}

		now := p.now().UTC()
		address.UpdatedAt = &now

		query := `
			UPDATE addresses SET
				uf = $1,
				localidade_numero = $2,
				localidade = $3,
				logradouro_dnec = $4,
				bairro_numero = $5,
				bairro = $6,
				faixas_caixa_postal = $7,
				faixas_cep = $8,
				loc_no_sem = $9,
				localidade_subordinada = $10,
				logradouro_texto = $11,
				logradouro_texto_adicional = $12,
				nome_unidade = $13,
				numero_localidade = $14,
				situacao = $15,
				tipo_cep = $16,
				desmembramento = $17,
				updated_at = $18
			WHERE
				cep = $19;
		`

		_, err = tx.Exec(ctx, query,
			address.UF,
			address.LocalidadeNumero,
			address.Localidade,
			address.LogradouroDNEC,
			address.BairoNumero,
			address.Bairro,
			address.FaixasCaixaPostal,
			address.FaixasCEP,
			address.LocNoSem,
			address.LocalidadeSubordinada,
			address.LogradouroTexto,
			address.LogradouroTextoAdicional,
			address.NomeUnidade,
			address.NumeroLocalidade,
			address.Situacao,
			address.TipoCEP,
			address.Desmembramento,
			address.UpdatedAt,
			address.UpdatedAt,
			cep,
		)
		if err != nil {
			return errors.E(op, kind(err), err)
		}

		return nil
	}

	return p.ExecTx(ctx, updateFn, op)
}

func (p *Postgres) Get(ctx context.Context, cep string) (*storage.Address, error) {
	const op errors.Op = "postgres.GetAddressCodeByHandle"
	return p.get(ctx, cep, op)
}

func (p *Postgres) get(ctx context.Context, cep string, op errors.Op) (*storage.Address, error) {
	query := `
		SELECT 
			cep,
			uf,
			localidade_numero,
			localidade,
			logradouro_dnec,
			bairro_numero,
			bairro,
			faixas_caixa_postal,
			faixas_cep,
			loc_no_sem,
			localidade_subordinada,
			logradouro_texto,
			logradouro_texto_adicional,
			nome_unidade,
			numero_localidade,
			situacao,
			tipo_cep,
			desmembramento,
			created_at,
			updated_at
		FROM addresses
		WHERE
			cep = $1;
	`
	row := p.db.QueryRow(ctx, query, cep)
	return scan(row, op)
}

func (p *Postgres) List(ctx context.Context, params storage.ListParams) ([]*storage.Address, error) {
	const op errors.Op = "postgres.ListAddressCodes"

	if params.Pagination == nil {
		params.Pagination = storage.NewPagination(storage.PaginationLimit, 0)
	}

	query := `
		SELECT 
			cep,
			uf,
			localidade_numero,
			localidade,
			logradouro_dnec,
			bairro_numero,
			bairro,
			faixas_caixa_postal,
			faixas_cep,
			loc_no_sem,
			localidade_subordinada,
			logradouro_texto,
			logradouro_texto_adicional,
			nome_unidade,
			numero_localidade,
			situacao,
			tipo_cep,
			desmembramento,
			created_at,
			updated_at
		FROM addresses WHERE uf = $1 ORDER BY cep ASC LIMIT $2 OFFSET $3;
	`
	p.db.Query(ctx, query, params.Pagination.Limit, params.Pagination.Offset)

	rows, err := p.db.Query(ctx, query, params.Pagination.Limit, params.Pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := make([]*storage.Address, 0)
	for rows.Next() {
		p, err := scan(rows, op)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, p)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.E(op, kind(err), err)
	}

	return addresses, nil
}

func scan(row pgx.Row, op errors.Op) (*storage.Address, error) {
	var address storage.Address
	err := row.Scan(
		&address.CEP,
		&address.UF,
		&address.LocalidadeNumero,
		&address.Localidade,
		&address.LogradouroDNEC,
		&address.BairoNumero,
		&address.Bairro,
		&address.FaixasCaixaPostal,
		&address.FaixasCEP,
		&address.LocNoSem,
		&address.LocalidadeSubordinada,
		&address.LogradouroTexto,
		&address.LogradouroTextoAdicional,
		&address.NomeUnidade,
		&address.NumeroLocalidade,
		&address.Situacao,
		&address.TipoCEP,
		&address.Desmembramento,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	if err != nil {
		return nil, errors.E(op, kind(err), err)
	}

	return &address, nil
}

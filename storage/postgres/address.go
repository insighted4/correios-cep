package postgres

import (
	"context"

	"github.com/insighted4/correios-cep/pkg/errors"
	"github.com/insighted4/correios-cep/storage"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) CreateAddress(ctx context.Context, address *storage.Address) error {
	const op errors.Op = "postgres.CreateAddress"
	query := `INSERT INTO addresses (
				cep,
				state,
				city,
				neighborhood,
            	location,
                children,
				created_at,
				updated_at
    		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8); 
	`

	now := p.now()
	address.CreatedAt = &now
	address.UpdatedAt = &now

	if _, err := p.db.Exec(ctx, query,
		address.CEP,
		address.State,
		address.City,
		address.Neighborhood,
		address.Location,
		address.Children,
		address.CreatedAt,
		address.UpdatedAt,
	); err != nil {
		return errors.E(op, kind(err), err)
	}

	return nil
}

func (p *Postgres) UpdateAddress(ctx context.Context, cep string, updater storage.Updater) error {
	const op errors.Op = "postgres.UpdateAddress"

	updateFn := func(tx pgx.Tx) error {
		old, err := p.get(ctx, cep, op)
		if err != nil {
			return err
		}

		address, err := updater(old)
		if err != nil {
			return err
		}

		now := p.now()
		address.UpdatedAt = &now

		query := `
			UPDATE addresses SET
				cep = $1,
				state = $2,
				city = $3,
				neighborhood = $4,
				location = $5,
				children = $6,
				updated_at = $7
			WHERE
				cep = $8;
		`

		_, err = tx.Exec(ctx, query,
			address.CEP,
			address.State,
			address.City,
			address.Neighborhood,
			address.Location,
			address.Children,
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

func (p *Postgres) GetAddress(ctx context.Context, cep string) (*storage.Address, error) {
	const op errors.Op = "postgres.GetAddress"
	return p.get(ctx, cep, op)
}

func (p *Postgres) get(ctx context.Context, cep string, op errors.Op) (*storage.Address, error) {
	query := `
		SELECT 
			cep,
			state,
			city,
			neighborhood,
			location,
			children,
			created_at,
			updated_at
		FROM addresses
		WHERE
			cep = $1;
	`
	row := p.db.QueryRow(ctx, query, cep)
	return scan(row, op)
}

func (p *Postgres) ListAddresses(ctx context.Context, params storage.ListParams) ([]*storage.Address, error) {
	const op errors.Op = "postgres.ListAddresses"

	if params.State == "" {
		return nil, errors.E(op, errors.KindBadRequest, "param state is required")
	}

	if params.Pagination == nil {
		return nil, errors.E(op, errors.KindUnexpected, "invalid pagination")
	}

	query := `
		SELECT 
			cep,
			state,
			city,
			neighborhood,
			location,
			children,
			created_at,
			updated_at
		FROM addresses WHERE lower(state) = lower($1) ORDER BY cep ASC LIMIT $2 OFFSET $3;
	`
	rows, err := p.db.Query(ctx, query, params.State, params.Pagination.Limit, params.Pagination.Offset)
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
	if err := row.Scan(
		&address.CEP,
		&address.State,
		&address.City,
		&address.Neighborhood,
		&address.Location,
		&address.Children,
		&address.CreatedAt,
		&address.UpdatedAt,
	); err != nil {
		return nil, errors.E(op, kind(err), err)
	}

	return &address, nil
}

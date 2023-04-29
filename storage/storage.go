package storage

import (
	"context"
)

type Storage interface {
	Close() error
	Check(ctx context.Context) error

	Create(ctx context.Context, code *Address) error
	Update(ctx context.Context, cep string, updater Updater) error
	Get(ctx context.Context, cep string) (*Address, error)
	List(ctx context.Context, params ListParams) ([]*Address, error)
}

type (
	Updater func(old *Address) (*Address, error)

	ListParams struct {
		Pagination *Pagination
		UF         string
	}
)

const (
	PaginationLimit = 100
)

// Pagination is passed as a parameter to limit the total of rows.
type Pagination struct {
	Limit  int
	Offset int
}

func NewPagination(perPage, page int) *Pagination {
	if perPage >= PaginationLimit {
		perPage = PaginationLimit
	}

	return &Pagination{
		Limit:  perPage,
		Offset: page * perPage,
	}
}

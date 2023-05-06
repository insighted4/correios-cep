package storage

import (
	"context"
)

type Storage interface {
	Close()
	Check(ctx context.Context) error

	CreateAddress(ctx context.Context, address *Address) error
	UpdateAddress(ctx context.Context, cep string, updater Updater) error
	GetAddress(ctx context.Context, cep string) (*Address, error)
	ListAddresses(ctx context.Context, params ListParams) ([]*Address, error)
}

type (
	Updater func(old *Address) (*Address, error)

	ListParams struct {
		State      string
		Pagination *Pagination
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
	if perPage < 1 || perPage > PaginationLimit {
		perPage = PaginationLimit
	}

	return &Pagination{
		Limit:  perPage,
		Offset: page * perPage,
	}
}

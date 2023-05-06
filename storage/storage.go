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

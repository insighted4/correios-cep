package correios

import (
	"context"

	"github.com/insighted4/correios-cep/storage"
)

type Correios interface {
	Check(ctx context.Context) error
	Lookup(ctx context.Context, cep string) (*storage.Address, error)
}

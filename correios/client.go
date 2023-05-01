package correios

import (
	"context"

	"github.com/insighted4/correios-cep/storage"
)

type Client interface {
	Check(ctx context.Context) error
	Lookup(ctx context.Context, cep string) (*storage.Address, error)
}

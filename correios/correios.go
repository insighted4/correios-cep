package correios

import (
	"context"
	"fmt"
	"net/http"

	"github.com/insighted4/correios-cep/pkg/errors"
	"github.com/insighted4/correios-cep/storage"
)

const (
	lookupURL = "https://buscacepinter.correios.com.br/app/consulta/html/consulta-detalhes-cep.php"
)

type Correios struct {
	cli *http.Client
}

var _ Client = (*Correios)(nil)

func New(cli *http.Client) *Correios {
	return &Correios{cli: cli}
}

func (c *Correios) Check(ctx context.Context) error {
	const op errors.Op = "correios.Check"

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, lookupURL, nil)
	if err != nil {
		return errors.E(op, errors.KindUnexpected, err)
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		return errors.E(op, errors.KindUnexpected, err)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.E(op, errors.KindUnexpected, fmt.Sprintf("HEAD %s returned %d", lookupURL, resp.StatusCode))
	}

	return nil
}

func (c *Correios) Lookup(ctx context.Context, cep string) (*storage.Address, error) {
	return nil, nil
}

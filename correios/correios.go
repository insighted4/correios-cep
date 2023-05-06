package correios

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/insighted4/correios-cep/pkg/errors"
	"github.com/insighted4/correios-cep/pkg/log"
	"github.com/insighted4/correios-cep/storage"
	"github.com/sirupsen/logrus"
)

const (
	lookupURL = "https://buscacepinter.correios.com.br/app/consulta/html/consulta-detalhes-cep.php"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:10.0) Gecko/20100101 Firefox/10.0"
)

type client struct {
	cli    *http.Client
	logger logrus.FieldLogger
}

var _ Correios = (*client)(nil)

func New(cli *http.Client) Correios {
	return &client{cli: cli, logger: log.Logger().WithField("component", "correios")}
}

func (c *client) Check(ctx context.Context) error {
	const op errors.Op = "correios.Check"

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, lookupURL, nil)
	if err != nil {
		c.logger.Errorf("unable to create new request: %v", err)
		return errors.E(op, errors.KindUnexpected, err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.cli.Do(req)
	if err != nil {
		c.logger.Errorf("failed to check Correios: %v", err)
		return errors.E(op, errors.KindUnexpected, err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Errorf("failed to check Correios: unexpected status code %d", resp.StatusCode)
		return errors.E(op, errors.KindUnexpected, fmt.Sprintf("HEAD %s returned %d", lookupURL, resp.StatusCode))
	}

	return nil
}

type Dado struct {
	UF                       string        `json:"uf"`
	Localidade               string        `json:"localidade"`
	LocNoSem                 string        `json:"locNoSem"`
	LocNu                    string        `json:"locNu"`
	LocalidadeSubordinada    string        `json:"localidadeSubordinada"`
	LogradouroDNEC           string        `json:"logradouroDNEC"`
	LogradouroTextoAdicional string        `json:"logradouroTextoAdicional"`
	LogradouroTexto          string        `json:"logradouroTexto"`
	Bairro                   string        `json:"bairro"`
	BaiNu                    string        `json:"baiNu"`
	NomeUnidade              string        `json:"nomeUnidade"`
	CEP                      string        `json:"cep"`
	TipoCEP                  string        `json:"tipoCep"`
	NumeroLocalidade         string        `json:"numeroLocalidade"`
	Situacao                 string        `json:"situacao"`
	FaixasCaixaPostal        []interface{} `json:"faixasCaixaPostal"`
	FaixasCEP                []interface{} `json:"faixasCep"`
}

type Response struct {
	Erro     bool    `json:"erro"`
	Mensagem string  `json:"mensagem"`
	Total    int     `json:"total"`
	Dados    []*Dado `json:"dados"`
}

func (c *client) Lookup(ctx context.Context, cep string) (*storage.Address, error) {
	const op errors.Op = "correios.Lookup"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, lookupURL, nil)
	if err != nil {
		c.logger.Errorf("unable to create new request: %v", err)
		return nil, errors.E(op, errors.KindUnexpected, err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")

	req.PostForm = url.Values{}
	req.PostForm.Set("cep", cep)

	resp, err := c.cli.Do(req)
	if err != nil {
		c.logger.Errorf("failed to lookup address: %v", err)
		return nil, errors.E(op, errors.KindUnexpected, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		c.logger.Errorf("failed to lookup address: unexpected status code %d", resp.StatusCode)
		return nil, errors.E(op, errors.KindUnexpected, fmt.Sprintf("POST %s returned %d", lookupURL, resp.StatusCode))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.E(op, errors.KindUnexpected)
	}

	response := new(Response)
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, errors.E(op, errors.KindUnexpected, err)
	}

	if len(response.Dados) == 0 {
		return nil, errors.E(op, errors.KindNotFound, fmt.Sprintf("cep %s not found", cep))
	}

	address := &storage.Address{CEP: cep}
	for _, dado := range response.Dados {
		address.Addresses = append(address.Addresses, &storage.Address{
			CEP:          dado.CEP,
			State:        dado.UF,
			City:         dado.Localidade,
			Neighborhood: dado.Bairro,
			Location:     dado.LogradouroDNEC,
		})
	}

	return address, nil
}

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

package correios

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/insighted4/correios-cep/pkg/errors"
	"github.com/insighted4/correios-cep/pkg/log"
	"github.com/insighted4/correios-cep/pkg/net"
	"github.com/insighted4/correios-cep/storage"
	"github.com/sirupsen/logrus"
)

const (
	baseURL   = "https://buscacepinter.correios.com.br"
	lookupURL = "/app/consulta/html/consulta-detalhes-cep.php"
)

type client struct {
	cli    *resty.Client
	logger logrus.FieldLogger
}

var _ Correios = (*client)(nil)

func New() Correios {
	logger := log.WithField("component", "correios")
	cli := net.NewClient().
		SetBaseURL(baseURL).
		SetHeader("Accept", "application/json").
		SetHeader("Referer", baseURL).
		SetLogger(logger).
		SetRetryCount(3).
		SetTimeout(20 * time.Second)

	return &client{
		cli:    cli,
		logger: logger,
	}
}

func (c *client) Check(ctx context.Context) error {
	const op errors.Op = "correios.Check"

	resp, err := c.cli.R().Head(lookupURL)
	if err != nil {
		c.logger.Errorf("failed to check Correios: %v", err)
		return errors.E(op, errors.KindUnexpected, err)
	}

	if resp.StatusCode() != http.StatusOK {
		c.logger.Errorf("failed to check Correios: unexpected status code %d", resp.StatusCode)
		return errors.E(op, errors.KindUnexpected, fmt.Sprintf("HEAD %s returned %d", lookupURL, resp.StatusCode()))
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

func (d *Dado) toAddress() *storage.Address {
	return &storage.Address{
		CEP:          d.CEP,
		State:        d.UF,
		City:         d.Localidade,
		Neighborhood: d.Bairro,
		Location:     d.LogradouroDNEC,
	}
}

type LookupResponse struct {
	Erro     bool    `json:"erro"`
	Mensagem string  `json:"mensagem"`
	Total    int     `json:"total"`
	Dados    []*Dado `json:"dados"`
}

func (c *client) Lookup(ctx context.Context, cep string) (*storage.Address, error) {
	const op errors.Op = "correios.Lookup"

	form := map[string]string{
		"cep": cep,
	}

	response := new(LookupResponse)
	resp, err := c.cli.R().SetFormData(form).SetResult(&response).Post(lookupURL)
	if err != nil {
		c.logger.Errorf("failed to lookup address: %v", err)
		return nil, errors.E(op, errors.KindUnexpected, err)
	}

	if resp.StatusCode() != http.StatusOK {
		c.logger.Errorf("failed to lookup address: unexpected status code %d", resp.StatusCode)
		return nil, errors.E(op, errors.KindUnexpected, fmt.Sprintf("POST %s returned %d", lookupURL, resp.StatusCode()))
	}

	if len(response.Dados) == 0 {
		return nil, errors.E(op, errors.KindNotFound, fmt.Sprintf("cep %s not found", cep))
	}

	address := &storage.Address{CEP: cep}
	if len(response.Dados) == 1 {
		return response.Dados[0].toAddress(), nil
	}

	for _, dado := range response.Dados {
		address.Children = append(address.Children, dado.toAddress())
	}

	return address, nil
}

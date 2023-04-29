package storage

import "time"

type Address struct {
	CEP                      string        `json:"cep,omitempty" db:"cep"`
	UF                       string        `json:"uf,omitempty"  db:"uf"`
	LocalidadeNumero         string        `json:"localidade_numero,omitempty"  db:"localidade_numero"`
	Localidade               string        `json:"localidade,omitempty"  db:"localidade"`
	LogradouroDNEC           string        `json:"logradouro_dnec,omitempty"  db:"logradouro_dnec"`
	BairoNumero              string        `json:"bairro_numero,omitempty"  db:"bairro_numero"`
	Bairro                   string        `json:"bairro,omitempty"  db:"bairro"`
	FaixasCaixaPostal        []interface{} `json:"faixas_caixa_postal,omitempty"  db:"faixas_caixa_postal"`
	FaixasCEP                []interface{} `json:"faixas_cep,omitempty"  db:"faixas_cep"`
	LocNoSem                 string        `json:"loc_no_sem,omitempty"  db:"loc_no_sem"`
	LocalidadeSubordinada    string        `json:"localidade_subordinada,omitempty"  db:"localidade_subordinada"`
	LogradouroTexto          string        `json:"logradouro_texto,omitempty"  db:"logradouro_texto"`
	LogradouroTextoAdicional string        `json:"logradouro_texto_adicional,omitempty"  db:"logradouro_texto_adicional"`
	NomeUnidade              string        `json:"nome_unidade,omitempty"  db:"nome_unidade"`
	NumeroLocalidade         string        `json:"numero_localidade,omitempty"  db:"numero_localidade"`
	Situacao                 string        `json:"situacao,omitempty"  db:"situacao"`
	TipoCEP                  string        `json:"tipo_cep,omitempty"  db:"tipo_cep"`
	Desmembramento           []interface{} `json:"desmembramento,omitempty"  db:"desmembramento"`
	CreatedAt                *time.Time    `json:"created_at,omitempty,omitempty"  db:"cep"`
	UpdatedAt                *time.Time    `json:"updated_at,omitempty,omitempty"  db:"cep"`
}

func (a Address) IsDesmembrado() bool {
	return len(a.Desmembramento) > 0
}

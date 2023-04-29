CREATE TABLE IF NOT EXISTS cages
(
    cep                        TEXT                      NOT NULL PRIMARY KEY,
    uf                         TEXT                      NOT NULL,
    localidade_numero          TEXT,
    localidade                 TEXT                      NOT NULL,
    logradouro_dnec            TEXT,
    bairro_numero              TEXT,
    bairro                     TEXT,
    faixas_caixa_postal        BYTEA,
    faixas_cep                 BYTEA,
    loc_no_sem                 TEXT,
    localidade_subordinada     TEXT,
    logradouro_texto           TEXT,
    logradouro_texto_adicional TEXT,
    nome_unidade               TEXT,
    numero_localidade          TEXT,
    situacao                   TEXT,
    tipo_cep                   TEXT,
    desmembramento             BYTEA,

    created_at                 TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at                 TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE INDEX IF NOT EXISTS addresses_uf_idx ON addresses (uf);

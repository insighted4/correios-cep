CREATE TABLE IF NOT EXISTS addresses
(
    cep          TEXT                      NOT NULL PRIMARY KEY,
    state        TEXT,
    city         TEXT,
    neighborhood TEXT,
    location     TEXT,
    source       JSONB,
    created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE INDEX IF NOT EXISTS addresses_uf_idx ON addresses (state);
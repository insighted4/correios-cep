package storage

import (
	"time"
)

type Address struct {
	CEP          string     `json:"cep" db:"cep"`
	State        string     `json:"state"  db:"state"`
	City         string     `json:"city" db:"location"`
	Neighborhood string     `json:"neighborhood" db:"neighborhood"`
	Location     string     `json:"location" db:"location"`
	Addresses    []*Address `json:"addresses" db:"addresses"`

	CreatedAt *time.Time `json:"created_at,omitempty,omitempty"  db:"cep"`
	UpdatedAt *time.Time `json:"updated_at,omitempty,omitempty"  db:"cep"`
}

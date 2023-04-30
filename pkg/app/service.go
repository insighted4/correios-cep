package app

import (
	"context"
	"time"
)

const Description = "Correios CEP Admin (API)"

type Checker interface {
	Check(ctx context.Context) error
}

func StartDate() time.Time {
	date, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00+00:00")
	if err != nil {
		panic(err)
	}

	return date.UTC()
}

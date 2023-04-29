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
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	return time.Date(2023, 01, 01, 01, 01, 0, 0, loc)
}

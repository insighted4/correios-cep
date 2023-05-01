package app

import (
	"time"
)

const Description = "Correios CEP Admin (API)"

func StartDate() time.Time {
	date, err := time.Parse(time.RFC3339, "2023-01-01T00:00:00+00:00")
	if err != nil {
		panic(err)
	}

	return date.UTC()
}

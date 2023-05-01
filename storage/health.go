package storage

import (
	"context"
	"time"
)

// NewCustomHealthCheckFunc returns a new health check function.
func NewCustomHealthCheckFunc(s Storage, now func() time.Time) func(context.Context) (details interface{}, err error) {
	return func(ctx context.Context) (details interface{}, err error) {
		if err := s.Check(ctx); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

package health

import (
	"context"
	"time"
)

type Checker interface {
	Check(ctx context.Context) error
}

// NewCustomHealthCheckFunc returns a new health check function.
func NewCustomHealthCheckFunc(checker Checker, now func() time.Time) func(context.Context) (details interface{}, err error) {
	return func(ctx context.Context) (details interface{}, err error) {
		if err := checker.Check(ctx); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

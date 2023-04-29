package stats

import (
	"context"

	"go.opencensus.io/trace"
)

// StartSpan takes in a Context Interface and opName and starts a span. It returns the new attached ObserverContext
// and span
func StartSpan(ctx context.Context, op string) (context.Context, *trace.Span) {
	return trace.StartSpan(ctx, op)
}

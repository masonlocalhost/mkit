package tracing

import "context"

const (
	skipTraceContextKey = "skip-tracing"
)

func WithSkipTrace(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipTraceContextKey, true)
}

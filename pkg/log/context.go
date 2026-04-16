package log

import (
	"context"
	"log/slog"
)

type loggerKey struct{}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// GetLogger retrieves the logger stored in ctx. Falls back to the default
// slog logger if none was stored.
func GetLogger(ctx context.Context) *slog.Logger {
	if v := ctx.Value(loggerKey{}); v != nil {
		return v.(*slog.Logger)
	}
	return slog.Default()
}

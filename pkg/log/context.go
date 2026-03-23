package log

import (
	"context"
	"github.com/sirupsen/logrus"
)

var l = &logrus.Entry{
	Logger: logrus.StandardLogger(),
	// Default is three fields plus a little extra room.
	Data: make(map[string]any, 6),
}

type loggerKey struct{}

func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger.WithContext(ctx))
}

// GetLogger retrieves the current logger from the context. If no logger is
// available, the default logger is returned.
func GetLogger(ctx context.Context) *logrus.Entry {
	if logger := ctx.Value(loggerKey{}); logger != nil {
		return logger.(*logrus.Entry)
	}
	return l.WithContext(ctx)
}

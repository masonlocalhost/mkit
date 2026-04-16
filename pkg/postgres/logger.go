package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormSlogLogger struct {
	LogLevel logger.LogLevel
	Logger   *slog.Logger
}

func (l *GormSlogLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *GormSlogLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Logger.InfoContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *GormSlogLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Logger.WarnContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *GormSlogLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Logger.ErrorContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *GormSlogLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel < logger.Info {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			l.Logger.ErrorContext(ctx, "gorm error",
				"duration", elapsed,
				"rows", rows,
				"sql", sql,
				"error", err,
			)
		}
		return
	}

	l.Logger.InfoContext(ctx, "gorm query",
		"duration", elapsed,
		"rows", rows,
		"sql", sql,
	)
}

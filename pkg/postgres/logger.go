package postgres

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type GormLogrusLogger struct {
	LogLevel logger.LogLevel
	Logger   *logrus.Logger
}

func (l *GormLogrusLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level

	return l
}

func (l *GormLogrusLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Logger.WithContext(ctx).Infof(msg, data...)
	}
}

func (l *GormLogrusLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Logger.WithContext(ctx).Warnf(msg, data...)
	}
}

func (l *GormLogrusLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Logger.WithContext(ctx).Errorf(msg, data...)
	}
}

func (l *GormLogrusLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	entry := l.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"duration": elapsed,
		"rows":     rows,
		"sql":      sql,
	})

	if l.LogLevel >= logger.Info {
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				entry.WithField("error", err).Info("gorm error")
			}
		} else {
			entry.Info("gorm query")
		}
	}
}

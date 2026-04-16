package rabbitmq

import (
	"fmt"
	"log/slog"
)

// Logger adapts *slog.Logger to the interface expected by go-rabbitmq
// (Fatalf / Errorf / Warnf / Infof / Debugf).
type Logger struct {
	logger   *slog.Logger
	minLevel slog.Level
}

func NewLogger(logger *slog.Logger, minLevel slog.Level) *Logger {
	return &Logger{logger: logger, minLevel: minLevel}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.minLevel <= slog.LevelError {
		l.logger.Error(fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.minLevel <= slog.LevelError {
		l.logger.Error(fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.minLevel <= slog.LevelWarn {
		l.logger.Warn(fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.minLevel <= slog.LevelInfo {
		l.logger.Info(fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.minLevel <= slog.LevelDebug {
		l.logger.Debug(fmt.Sprintf(format, v...))
	}
}

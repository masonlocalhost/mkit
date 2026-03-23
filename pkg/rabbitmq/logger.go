package rabbitmq

import "github.com/sirupsen/logrus"

type Logger struct {
	Logger *logrus.Logger
	Level  logrus.Level
}

func NewLogger(logger *logrus.Logger, level logrus.Level) *Logger {
	return &Logger{
		Logger: logger,
		Level:  level,
	}
}

func (l Logger) Fatalf(format string, v ...interface{}) {
	if l.Level >= logrus.FatalLevel {
		l.Logger.Fatalf(format, v...)
	}
}

func (l Logger) Errorf(format string, v ...interface{}) {
	if l.Level >= logrus.ErrorLevel {
		l.Logger.Errorf(format, v...)
	}
}

func (l Logger) Warnf(format string, v ...interface{}) {
	if l.Level >= logrus.WarnLevel {
		l.Logger.Warnf(format, v...)
	}
}

func (l Logger) Infof(format string, v ...interface{}) {
	if l.Level >= logrus.InfoLevel {
		l.Logger.Infof(format, v...)
	}
}

func (l Logger) Debugf(format string, v ...interface{}) {
	if l.Level >= logrus.DebugLevel {
		l.Logger.Debugf(format, v...)
	}
}

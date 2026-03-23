package log

import (
	"fmt"
	"io"
	"mkit/pkg/config"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// NewLogger create an logrus instance for using globally
func NewLogger(cfg *config.App) (*logrus.Logger, error) {
	var (
		logCfg  = cfg.Log
		logger  = logrus.New()
		logPath = logCfg.LogFilePath
	)

	level, err := logrus.ParseLevel(logCfg.Level)
	if err != nil {
		return nil, fmt.Errorf("cant parse logrus level for '%s': %w", level, err)
	}

	if err := os.MkdirAll(logPath, 0777); err != nil {
		return nil, fmt.Errorf("cannot create log directory: %w", err)
	}

	logger.SetLevel(level)
	logger.SetFormatter(&UTCFormatter{Formatter: NewDetailFormatter("")})
	logger.SetReportCaller(true)

	writer, err := rotatelogs.New(
		fmt.Sprintf("%v/%v", logPath, "log-%Y-%m-%d.txt"),
		rotatelogs.WithMaxAge(time.Hour*24*31),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithClock(rotatelogs.UTC),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create rotate log writer: %w", err)
	}

	logger.SetOutput(io.MultiWriter(writer, os.Stdout))

	return logger, nil
}

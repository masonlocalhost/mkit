package log

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	otellog "go.opentelemetry.io/otel/sdk/log"
	"mkit/pkg/config"
)

// NewLogger creates a *slog.Logger that:
//   - writes colourised, human-readable output to stdout (via tint)
//   - forwards every record to OpenTelemetry as raw JSON over OTLP when
//     logProvider is non-nil (set up by the tracing package)
func NewLogger(cfg *config.App, logProvider *otellog.LoggerProvider) (*slog.Logger, error) {
	level, err := parseLevel(cfg.Log.Level)
	if err != nil {
		return nil, fmt.Errorf("cannot parse log level %q: %w", cfg.Log.Level, err)
	}

	termHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      level,
		TimeFormat: "Jan _2 15:04:05.000",
		AddSource:  true,
	})

	handlers := []slog.Handler{termHandler}

	if logProvider != nil {
		otelHandler := otelslog.NewHandler(
			cfg.Tracing.ServiceName,
			otelslog.WithLoggerProvider(logProvider),
			otelslog.WithVersion(cfg.Version),
		)
		handlers = append(handlers, otelHandler)
	}

	logger := slog.New(newMultiHandler(handlers...))
	slog.SetDefault(logger)

	return logger, nil
}

// parseLevel maps logrus-compatible level strings to slog.Level.
func parseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(s) {
	case "trace", "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error", "fatal", "panic":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unknown level %q", s)
	}
}

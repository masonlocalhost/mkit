package chi

import (
	"log/slog"
	"mkit/pkg/config"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func New(cfg *config.App, logger *slog.Logger) (chi.Router, error) {
	r := chi.NewRouter()

	r.Use(CORS())
	if cfg.Tracing.Enabled {
		r.Use(func(next http.Handler) http.Handler {
			return otelhttp.NewHandler(next, cfg.Tracing.ServiceName)
		})
	}
	r.Use(LoggerMiddleware(logger))
	r.Use(RecoveryMiddleware(logger))

	return r, nil
}

package connect

import (
	"fmt"
	"log/slog"
	"mkit/pkg/config"
	"net/http"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func New(cfg *config.App, logger *slog.Logger) (*http.ServeMux, []connect.HandlerOption, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, nil, fmt.Errorf("cant setup validator: %w", err)
	}

	opts := []connect.HandlerOption{
		connect.WithInterceptors(
			UnaryPanicInterceptor(logger),
			UnaryLogger(logger),
			UnaryValidation(validator),
		),
	}

	return http.NewServeMux(), opts, nil
}

// WrapMux wraps the mux with OpenTelemetry HTTP instrumentation when tracing is enabled.
func WrapMux(mux *http.ServeMux, cfg *config.App) http.Handler {
	if cfg.Tracing.Enabled {
		return otelhttp.NewHandler(mux, cfg.Tracing.ServiceName)
	}
	return mux
}

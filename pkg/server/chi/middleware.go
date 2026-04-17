package chi

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"mkit/pkg/log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
)

func LoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				id, _ := uuid.NewV7()
				reqID = id.String()
			}

			entry := logger.With("request_id", reqID)
			ctx := log.WithLogger(r.Context(), entry)
			r = r.WithContext(ctx)

			ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()
			next.ServeHTTP(ww, r)

			entry.InfoContext(r.Context(), "Incoming request",
				"status", ww.status,
				"duration", time.Since(start),
				"method", r.Method,
				"path", r.URL.Path,
				"client_ip", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		})
	}
}

func RecoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					stack := string(debug.Stack())
					logger.ErrorContext(r.Context(), fmt.Sprintf("panic recovered: %v", rec),
						"stack", stack,
						"url", r.URL.String(),
					)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(map[string]any{
						"error":  "Internal Server Error",
						"errors": []string{"server panicked"},
					})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, access_token, token, X-Requested-With")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Length")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

package gin

import (
	"fmt"
	"log/slog"
	"mkit/pkg/log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			id, _ := uuid.NewV7()
			reqID = id.String()
		}

		entry := logger.With("request_id", reqID)
		ctx := log.WithLogger(c.Request.Context(), entry)
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()
		c.Next()

		entry.InfoContext(c.Request.Context(), "Incoming request",
			"status", c.Writer.Status(),
			"duration", time.Since(start),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)
	}
}

func RecoveryMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := string(debug.Stack())

				logger.ErrorContext(c.Request.Context(), fmt.Sprintf("panic recovered in Gin handler: %v", rec),
					"stack", stack,
					"url", c.Request.URL.String(),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":  "Internal Server Error",
					"errors": []string{"server panicked"},
				})
			}
		}()

		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	var cfg = cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
			http.MethodOptions, http.MethodPatch, http.MethodHead,
		},
		AllowHeaders: []string{
			"Origin", "Authorization", "Content-Type", "access_token",
			"token", "X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}

	return cors.New(cfg)
}

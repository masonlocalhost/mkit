package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.InfoContext(c.Request.Context(), "Incoming request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"duration", time.Since(start),
			"user_agent", c.Request.UserAgent(),
		)
	}
}

func Recovery(logger *slog.Logger) gin.HandlerFunc {
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
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type", "access_token",
			"token", "X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}

	return cors.New(cfg)
}

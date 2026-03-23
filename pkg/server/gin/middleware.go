package gin

import (
	"fmt"
	"mkit/pkg/log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			id, _ := uuid.NewV7()
			reqID = id.String()
		}

		start := time.Now()
		entry := logrus.NewEntry(logger).WithFields(logrus.Fields{
			"request_id": reqID,
		})

		ctx := log.WithLogger(c.Request.Context(), entry)
		c.Request = c.Request.WithContext(ctx)
		// Process request
		c.Next()

		duration := time.Since(start)
		entry.WithContext(c.Request.Context()).WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"duration":   duration,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
		}).Info("Incoming request")
	}
}

func RecoveryMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				// Stack trace
				stack := string(debug.Stack())

				logger.WithContext(c.Request.Context()).WithFields(logrus.Fields{
					"stack": stack,
					"url":   c.Request.URL.String(),
				}).Error(fmt.Sprintf("panic recovered in Gin handler: %v", rec))

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

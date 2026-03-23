package middleware

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
	"time"
)

func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		duration := time.Since(start)
		logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"duration":   duration,
			"user-agent": c.Request.UserAgent(),
		}).Info("Incoming request")
	}
}

func Recovery(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				// Stack trace
				stack := string(debug.Stack())

				logger.WithFields(logrus.Fields{
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

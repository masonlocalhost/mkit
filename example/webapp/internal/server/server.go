package server

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	logger      *slog.Logger
	redisClient *redis.Client
	db          *gorm.DB
}

func (s *Server) Name() string {
	return "webapp"
}

func (s *Server) RegisterRouter(router *gin.Engine) {
	// register routes, can be attached to Server as receiver functions or decoupled to controllers
	// simple example:
	router.GET("api/v1/ping", s.Ping)
}

func (s *Server) Init() error {
	s.logger.Info("Init server")
	return nil
}

func (s *Server) Close() error {
	s.logger.Info("Close server")
	return nil
}

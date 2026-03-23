package server

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server struct {
	logger      *logrus.Logger
	redisClient *redis.Client
	db          *gorm.DB
}

func (s *Server) Name() string {
	return "webapp"
}

func (s *Server) RegisterRouter(router *gin.Engine) {
	// register routes, can be attached to Server as receiver functions or decouped to controllers
	// simple example:
	router.GET("api/v1/ping", s.Ping)
}

func (s *Server) Init() error {
	// initialization
	s.logger.Infof("Init server")

	return nil
}

func (s *Server) Close() error {
	// close server dependency if needed
	s.logger.Info("Close server")

	return nil
}

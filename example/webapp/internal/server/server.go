package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (s *Server) RegisterRouter(router chi.Router) {
	router.Get("/api/v1/ping", s.Ping)
}

func (s *Server) Init() error {
	s.logger.Info("Init server")
	return nil
}

func (s *Server) Close() error {
	s.logger.Info("Close server")
	return nil
}

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`"pong"`))
}

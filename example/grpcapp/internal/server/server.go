package server

import (
	"log/slog"
	"mkit/example/grpcapp/config"
	"mkit/example/grpcapp/internal/service/technology"
	"mkit/example/grpcapp/pkg/api/go/nanoid/v1"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	db                *gorm.DB
	logger            *slog.Logger
	cfg               *config.Config
	technologyService *technology.Service

	nanoid.UnimplementedNanoidServiceServer
}

func (s *Server) Name() string {
	return "nanoid"
}

func (s *Server) RegisterPBs(server *grpc.Server) {
	nanoid.RegisterNanoidServiceServer(server, s)
}

func (s *Server) Init() error {
	return nil
}

func (s *Server) Close() error {
	return nil
}

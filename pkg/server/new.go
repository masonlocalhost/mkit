package server

import (
	"connectrpc.com/vanguard/vanguardgrpc"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type internalServer interface {
	Name() string
	Close() error
	Init() error
}

type internalGRPCServer interface {
	internalServer
	RegisterPBs(server *grpc.Server)
}

type internalGinServer interface {
	internalServer
	RegisterRouter(router *gin.Engine)
}
type Server struct {
	GRPCNetListener      net.Listener
	GRPCTranscodeServer  *http.Server
	GRPCTranscodeHandler http.Handler
	HTTPServer           *http.Server

	Deps                *Dependencies
	internalGRPCServers []internalGRPCServer
	internalGINServer   internalGinServer
}

func NewServer(deps ...Dependency) *Server {
	d := NewDependencies()

	for _, dep := range deps {
		dep(d)
	}

	return &Server{
		Deps: d,
	}
}

// RegisterInternalGRPCServers to be used after all other pb.Register...Server (before server.Serve())
func (s *Server) RegisterInternalGRPCServers(iss ...internalGRPCServer) {
	for _, i := range iss {
		if i != nil {
			i.RegisterPBs(s.Deps.GRPCServer)
			s.internalGRPCServers = append(s.internalGRPCServers, i)
		}
	}

	grpcConfig := s.Deps.AppConfig.GRPC
	logger := s.Deps.Logger
	if grpcConfig != nil && grpcConfig.JsonTranscodeEnabled {
		transcoder, err := vanguardgrpc.NewTranscoder(s.Deps.GRPCServer)
		if err != nil {
			logger.Fatalf("cannot init gRPC trancoder: %v", err)
		}

		s.GRPCTranscodeHandler = h2c.NewHandler(transcoder, &http2.Server{})
	}
}

func (s *Server) RegisterInternalGinServer(server internalGinServer) {
	server.RegisterRouter(s.Deps.GinEngine)
	s.internalGINServer = server
}

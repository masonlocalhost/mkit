package server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

func (s *Server) Serve() {
	cfg := s.Deps.AppConfig
	s.init()
	if cfg.HTTP != nil {
		s.serveHTTP()
	}
	if cfg.GRPC != nil {
		s.serveGRPC()
	}
}

func (s *Server) init() {
	var logger = s.Deps.Logger

	if len(s.internalGRPCServers) > 0 {
		for _, is := range s.internalGRPCServers {
			if err := is.Init(); err != nil {
				logger.Fatalf("Error when init grpc service %s: %v", is.Name(), err)
			}
		}
	}

	if s.internalGINServer != nil {
		if err := s.internalGINServer.Init(); err != nil {
			logger.Fatalf("Error when init gin service %s: %v", s.internalGINServer.Name(), err)
		}
	}
}

func (s *Server) serveGRPC() {
	var (
		deps   = s.Deps
		logger = deps.Logger
		cfg    = deps.AppConfig.GRPC
	)

	if len(s.internalGRPCServers) == 0 {
		logger.Fatalf("Servers must be registered to be able to serve (use RegisterInternalGRPCServers())")
	}

	grpcAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	lis, err := net.Listen("tcp", grpcAddr)
	if lis == nil {
		logger.Fatal(fmt.Sprintf("Failed to listen to %s (it seems like the port is occupied)", grpcAddr))
	}
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to listen to %s: %v", grpcAddr, err))
	}

	s.GRPCNetListener = lis
	if !cfg.JsonTranscodeEnabled {
		go func() {
			logger.Infof("GRPC server running on %s", grpcAddr)
			if err := deps.GRPCServer.Serve(lis); err != nil {
				logger.Fatalf("Error when starting gRPC: %v", err)
			}
		}()

		return
	}

	httpSrv := http.Server{
		Handler: s.GRPCTranscodeHandler,
	}

	logger.Infof("GRPC server with transcoder running on %s", grpcAddr)
	go func() {
		if err := httpSrv.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Failed to serve with http trancoder: %v", err)
		}
	}()
}

func (s *Server) serveHTTP() {
	var (
		deps   = s.Deps
		logger = deps.Logger
		cfg    = deps.AppConfig.HTTP
	)

	if s.internalGINServer == nil {
		logger.Fatalf("Gin server must be registered to be able to serve (use RegisterInternalGinServer())")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("/", deps.GinEngine.Handler())

	httpServer := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	}
	s.HTTPServer = httpServer

	logger.Infof("HTTP server running on %s", httpServer.Addr)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("HTTP server failed to listen: %v", err)
		}
	}()
}

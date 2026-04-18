package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
				logger.Error("Error when init grpc service", "name", is.Name(), "error", err)
				os.Exit(1)
			}
		}
	}

	if s.internalHTTPServer != nil {
		if err := s.internalHTTPServer.Init(); err != nil {
			logger.Error("Error when init http service", "name", s.internalHTTPServer.Name(), "error", err)
			os.Exit(1)
		}
	}

	for _, is := range s.internalConnectServers {
		if err := is.Init(); err != nil {
			logger.Error("Error when init connect service", "name", is.Name(), "error", err)
			os.Exit(1)
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
		logger.Error("Servers must be registered to be able to serve (use RegisterInternalGRPCServers())")
		os.Exit(1)
	}

	grpcAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	lis, err := net.Listen("tcp", grpcAddr)
	if lis == nil {
		logger.Error("Failed to listen — port may be occupied", "addr", grpcAddr)
		os.Exit(1)
	}
	if err != nil {
		logger.Error("Failed to listen", "addr", grpcAddr, "error", err)
		os.Exit(1)
	}

	s.GRPCNetListener = lis
	if !cfg.JsonTranscodeEnabled {
		go func() {
			logger.Info("GRPC server running", "addr", grpcAddr)
			if err := deps.GRPCServer.Serve(lis); err != nil {
				logger.Error("Error when starting gRPC", "error", err)
				os.Exit(1)
			}
		}()

		return
	}

	httpSrv := http.Server{
		Handler: s.GRPCTranscodeHandler,
	}

	logger.Info("GRPC server with transcoder running", "addr", grpcAddr)
	go func() {
		if err := httpSrv.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Failed to serve with http transcoder", "error", err)
			os.Exit(1)
		}
	}()
}

func (s *Server) serveHTTP() {
	var (
		deps   = s.Deps
		logger = deps.Logger
		cfg    = deps.AppConfig.HTTP
	)

	hasConnect := len(s.internalConnectServers) > 0

	if s.internalHTTPServer == nil && !hasConnect {
		logger.Error("HTTP server must be registered to be able to serve (use RegisterInternalHTTPServer())")
		os.Exit(1)
	}

	// When ConnectRPC servers are registered, use their mux as the base so
	// their specific path patterns are matched before chi's catch-all "/".
	var mux *http.ServeMux
	if hasConnect {
		mux = deps.ConnectMux
	} else {
		mux = http.NewServeMux()
	}

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	if s.internalHTTPServer != nil {
		mux.Handle("/", deps.ChiRouter)
	}

	// h2c enables HTTP/2 cleartext, required for ConnectRPC's gRPC protocol.
	var handler http.Handler = mux
	if hasConnect {
		handler = h2c.NewHandler(mux, &http2.Server{})
	}

	var idleTimeout time.Duration
	if hasConnect && cfg.Connect != nil && cfg.Connect.MaxConnectionAge != "" {
		if d, err := time.ParseDuration(cfg.Connect.MaxConnectionAge); err == nil {
			idleTimeout = d
		}
	}

	httpServer := &http.Server{
		Handler:     handler,
		Addr:        fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		IdleTimeout: idleTimeout,
	}
	s.HTTPServer = httpServer

	logger.Info("HTTP server running", "addr", httpServer.Addr)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failed to listen", "error", err)
			os.Exit(1)
		}
	}()
}

// fatal logs at Error level then exits. Use only for unrecoverable startup errors.
func fatal(logger *slog.Logger, msg string, args ...any) {
	logger.Error(msg, args...)
	os.Exit(1)
}

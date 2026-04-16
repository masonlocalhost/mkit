package server

import (
	"context"
	stdlog "log"
	"os"

	"mkit/example/grpcapp/config"
	techrepo "mkit/example/grpcapp/internal/repository/technology"
	"mkit/example/grpcapp/internal/service/technology"

	"mkit/pkg/log"
	"mkit/pkg/postgres"
	"mkit/pkg/server"
	"mkit/pkg/server/grpc"
	"mkit/pkg/tracing"

	"os/signal"
	"syscall"
)

func Run() {
	var (
		ctx, stop = signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	)
	defer stop()

	cfg, err := config.GetConfig()
	if err != nil {
		stdlog.Fatalf("cannot init config: %v", err)
	}

	appCfg := &cfg.App

	// Init tracing first so the log provider is available for the logger.
	trace, err := tracing.NewService(ctx, appCfg)
	if err != nil {
		stdlog.Fatalf("failed to init tracing service: %v", err)
	}

	logger, err := log.NewLogger(appCfg, trace.LogProvider())
	if err != nil {
		stdlog.Fatalf("failed to init logger: %v", err)
	}

	db, err := postgres.New(logger, appCfg)
	if err != nil {
		logger.Error("failed to init postgresql db", "error", err)
		os.Exit(1)
	}

	grpcServer, healthServer, err := grpc.New(appCfg, logger)
	if err != nil {
		logger.Error("failed to init grpc server", "error", err)
		os.Exit(1)
	}

	service := server.NewServer(
		server.AppConfig(appCfg),
		server.Logger(logger),
		server.GRPCServer(grpcServer),
		server.Postgres(db),
		server.Tracing(trace),
		server.HealthServer(healthServer),
	)

	technologyRepo := techrepo.NewRepository(db)
	techService := technology.NewService(logger, technologyRepo)

	nanoidServer := &Server{
		db:                db,
		logger:            logger,
		cfg:               cfg,
		technologyService: techService,
	}

	service.RegisterInternalGRPCServers(nanoidServer)
	service.Serve()

	<-ctx.Done()
	service.Close(ctx)
}

package server

import (
	"context"
	stdlog "log"
	"os"

	"mkit/example/webapp/config"
	"mkit/pkg/cache/redis"
	"mkit/pkg/log"
	"mkit/pkg/postgres"
	"mkit/pkg/server"
	chi2 "mkit/pkg/server/chi"

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

	logger, err := log.NewLogger(appCfg, nil)
	if err != nil {
		stdlog.Fatalf("failed to init logger: %v", err)
	}

	db, err := postgres.New(logger, appCfg)
	if err != nil {
		logger.Error("failed to init postgresql db", "error", err)
		os.Exit(1)
	}

	redisClient, err := redis.NewClient(ctx, appCfg)
	if err != nil {
		logger.Error("failed to init redis", "error", err)
		os.Exit(1)
	}

	chiRouter, err := chi2.New(appCfg, logger)
	if err != nil {
		logger.Error("failed to init chi router", "error", err)
		os.Exit(1)
	}

	// mkit service handles lifecycle of common dependencies like: db, redis, tracing...
	service := server.NewServer(
		server.AppConfig(appCfg),
		server.Logger(logger),
		server.ChiRouter(chiRouter),
		server.Redis(redisClient),
		server.Postgres(db),
		// server.Tracing(trace),
	)

	// Setup other dependency (services, repositories), create server and register to service
	webAppServer := &Server{
		logger:      logger,
		redisClient: redisClient,
		db:          db,
		// other deps
	}

	service.RegisterInternalHTTPServer(webAppServer)
	service.Serve()

	<-ctx.Done()
	service.Close(ctx)
}

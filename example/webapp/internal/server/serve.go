package server

import (
	"context"
	slog "log"
	"mkit/example/webapp/config"
	"mkit/pkg/cache/redis"
	"mkit/pkg/log"
	"mkit/pkg/postgres"
	"mkit/pkg/server"
	"mkit/pkg/server/gin"
	"os"
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
		slog.Fatalf("cannot init config: %v", err)
	}

	appCfg := &cfg.App
	logger, err := log.NewLogger(appCfg)
	if err != nil {
		slog.Fatalf("failed to init logger: %v", err)
	}

	db, err := postgres.New(logger, appCfg)
	if err != nil {
		logger.Fatalf("failed to init postgresql db: %v", err)
	}

	redisClient, err := redis.NewClient(ctx, appCfg)
	if err != nil {
		logger.Fatalf("failed to init redis: %v", err)
	}

	ginEngine, err := gin.New(appCfg, logger)
	if err != nil {
		logger.Fatalf("failed to init gin engine db: %v", err)
	}

	// mkit service handles lifecycle of common depecdencies like: db, redis, tracing...
	service := server.NewServer(
		server.AppConfig(appCfg),
		server.Logger(logger),
		server.GinEngine(ginEngine),
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

	service.RegisterInternalGinServer(webAppServer)
	service.Serve()

	<-ctx.Done()
	service.Close(ctx)
}

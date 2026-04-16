package server

import (
	"context"
	stdlog "log"
	"os"

	"mkit/example/ginapp/config"
	ctrl "mkit/example/ginapp/internal/controller"
	"mkit/example/ginapp/internal/service/technology"

	// "mkit/pkg/cache/redis"
	// "mkit/pkg/cache/redsync"
	// "mkit/pkg/cron"
	"mkit/pkg/log"
	// "mkit/pkg/minio"
	techrepo "mkit/example/ginapp/internal/repository/technology"
	"mkit/pkg/postgres"
	gin2 "mkit/pkg/server/gin"

	// "mkit/pkg/pubsub"
	// "mkit/pkg/rabbitmq"
	"mkit/pkg/server"
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

	// migrate database
	// if err := postgres.MigrateDatabase(appCfg, logger); err != nil {
	// 	logger.Error("failed to migrate database", "error", err); os.Exit(1)
	// }

	// redisClient, err := redis.NewClient(ctx, appCfg)
	// if err != nil {
	// 	logger.Error("failed to init redis", "error", err); os.Exit(1)
	// }

	// rabbitMQ, err := rabbitmq.New(appCfg, logger)
	// if err != nil {
	// 	logger.Error("failed to init rabbitmq", "error", err); os.Exit(1)
	// }

	// pubSub, err := pubsub.NewService(
	// 	ctx, rabbitMQ, "nitro-dashboard", event.PubSubExchangeName, logger, event.ProtoRegistry,
	// )
	// if err != nil {
	// 	logger.Error("failed to init pubsub service", "error", err); os.Exit(1)
	// }

	db, err := postgres.New(logger, appCfg)
	if err != nil {
		logger.Error("failed to init postgresql db", "error", err)
		os.Exit(1)
	}

	ginEngine, err := gin2.New(appCfg, logger)
	if err != nil {
		logger.Error("failed to init gin engine", "error", err)
		os.Exit(1)
	}

	// loc, err := time.LoadLocation(cfg.Timezone)
	// if err != nil {
	// 	logger.Error("failed to load location", "error", err); os.Exit(1)
	// }

	// minioService, err := minio.New(appCfg, logger)
	// redsyncClient := redsync.NewClient(redisClient)

	// cronManager := cron.New(ctx, loc, logger, redsyncClient)

	service := server.NewServer(
		server.AppConfig(appCfg),
		server.Logger(logger),
		server.GinEngine(ginEngine),
		// server.Redis(redisClient),
		server.Postgres(db),
		server.Tracing(trace),
		// server.RabbitMQ(rabbitMQ),
		// server.CronManager(cronManager),
	)

	technologyRepo := techrepo.NewRepository(db)
	techService := technology.NewService(logger, technologyRepo)

	dep := &ctrl.DependencyContainer{
		Cfg:               cfg,
		Logger:            logger,
		DB:                db,
		TechnologyService: techService,
	}

	ginServer := &Server{dep: dep}

	service.RegisterInternalGinServer(ginServer)
	service.Serve()

	<-ctx.Done()
	service.Close(ctx)
}

package server

import (
	"context"
	slog "log"

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

	// migrate database
	// if err := postgres.MigrateDatabase(appCfg, logger); err != nil {
	// 	logger.Fatalf("failed to migrate database: %v", err)
	// }

	// redisClient, err := redis.NewClient(ctx, appCfg)
	// if err != nil {
	// 	logger.Fatalf("failed to init redis: %v", err)
	// }

	// rabbitMQ, err := rabbitmq.New(appCfg, logger)
	// if err != nil {
	// 	logger.Fatalf("failed to init rabbitmq: %v", err)
	// }

	// pubSub, err := pubsub.NewService(
	// 	ctx, rabbitMQ, "nitro-dashboard", event.PubSubExchangeName, logger, event.ProtoRegistry,
	// )
	// if err != nil {
	// 	logger.Fatalf("failed to init pubsub service: %v", err)
	// }

	db, err := postgres.New(logger, appCfg)
	if err != nil {
		logger.Fatalf("failed to init postgresql db: %v", err)
	}

	ginEngine, err := gin2.New(appCfg, logger)
	if err != nil {
		logger.Fatalf("failed to init gin engine db: %v", err)
	}

	trace, err := tracing.NewService(ctx, appCfg, logger)
	if err != nil {
		logger.Fatalf("failed to init tracing service: %v", err)
	}

	// loc, err := time.LoadLocation(cfg.Timezone)
	// if err != nil {
	// 	logger.Fatalf("failed to load location: %v", err)
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

	// userConfigRepo := userconfig.NewRepository(db)
	// userHistoryRepo := userhistory.NewRepository(db)
	// userInvokeTokenRepo := usertoken.NewRepository(db)
	// blacklistRepo := blacklistrepo.NewRepository(db)
	// queueItemRepo := queuerepo.NewRepository(db)
	// entityRepo := entityrepo.NewRepository(db)
	// collectionRepo := collectionrepo.NewRepository(db)
	technologyRepo := techrepo.NewRepository(db)
	// issueRepo := issuerepo.NewRepository(db)
	// issueResourceValidationRepo := resvalidationrepo.NewRepository(db)
	// hubService := hub.NewHub()
	// historyRepo := historyrepo.NewRepository(db)

	techService := technology.NewService(logger, technologyRepo)

	// maxmindCli, err := maxmindcli.New()
	// entityService := entity.NewService(logger, entityRepo, issueRepo, maxmindCli, collectionRepo, blacklistRepo, redisClient)
	// if err != nil {
	// 	logger.Fatalf("cant init maxmind cli: %v", err)
	// }
	// issueService := issue.NewService(logger, issueRepo, entityRepo, historyRepo)
	// if err != nil {
	// 	logger.Fatalf("cant init swagger ui service: %v", err)
	// }
	// issueResourceValidationService, err := resvalidation.NewService(logger, issueResourceValidationRepo, rabbitMQ, issueRepo)
	// if err != nil {
	// 	logger.Fatalf("failed to init issue resource validation service: %v", err)
	// }
	// historyService, err := history.NewService(ctx, logger, queueItemRepo, historyRepo, entityRepo, collectionRepo, redsyncClient)
	// if err != nil {
	// 	logger.Fatalf("cant init history service: %v", err)
	// }
	// casEzEngine, err := casez.NewClient(cfg)
	// if err != nil {
	// 	logger.Fatalf("cant init history service: %v", err)
	// }
	// swaggerUIService, err := swaggerui.NewService()
	// if err != nil {
	// 	logger.Fatalf("cant init swagger ui service: %v", err)
	// }
	// userService := user.NewService(
	// 	logger, userInvokeTokenRepo, userConfigRepo, userHistoryRepo, casEzEngine,
	// 	cfg.CasEzEngine.ServiceName, hubService, redisClient, collectionRepo,
	// )
	// authService := auth.NewService(
	// 	logger, userInvokeTokenRepo, userConfigRepo, userHistoryRepo, casEzEngine,
	// 	cfg.CasEzEngine.ServiceName, hubService,
	// )
	// orgService := org.NewService(logger, casEzEngine, cfg.CasEzEngine.ServiceName)
	// sseService := sse.NewService(logger)
	// collectionService, err := collection.NewService(
	// 	ctx, loc, logger, collectionRepo, blacklistRepo, redisClient, queueItemRepo, historyRepo,
	// 	userService, issueRepo, redsyncClient, pubSub,
	// )
	// if err != nil {
	// 	logger.Fatalf("cant init collection service: %v", err)
	// }

	dep := &ctrl.DependencyContainer{
		Cfg:    cfg,
		Logger: logger,
		// UserService:                    userService,
		// OrgService:                     orgService,
		// DiscoveryHistoryService:        historyService,
		// EntityService:                  entityService,
		// AuthService:                    authService,
		// Hub:                            hubService,
		DB: db,
		// CollectionService:              collectionService,
		// MaxmindCli:                     maxmindCli,
		// RedisClient:                    redisClient,
		TechnologyService: techService,
		// Minio:                          minioService,
		// IssueService:                   issueService,
		// SwaggerUIService:               swaggerUIService,
		// IssueResourceValidationService: issueResourceValidationService,
		// SSEService:                     sseService,
		// PubSub:                         pubSub,
		// CronManager:                    cronManager,
	}

	ginServer := &Server{dep: dep}

	service.RegisterInternalGinServer(ginServer)
	service.Serve()

	<-ctx.Done()
	service.Close(ctx)
}

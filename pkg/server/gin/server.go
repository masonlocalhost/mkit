package gin

import (
	"mkit/pkg/config"
	"mkit/pkg/enum"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func New(
	cfg *config.App, logger *logrus.Logger,
) (*gin.Engine, error) {
	if cfg.Environment == enum.EnvironmentProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	var (
		router   = gin.New()
		traceCfg = cfg.Tracing
	)

	router.Use(CORS())
	if traceCfg.Enabled {
		router.Use(otelgin.Middleware(traceCfg.ServiceName))
	}
	router.Use(LoggerMiddleware(logger))
	router.Use(RecoveryMiddleware(logger))

	return router, nil
}

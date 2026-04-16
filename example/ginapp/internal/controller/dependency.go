package controller

import (
	"log/slog"
	"mkit/example/ginapp/config"
	"mkit/example/ginapp/internal/service/technology"

	"gorm.io/gorm"
)

type DependencyContainer struct {
	Cfg               *config.Config
	Logger            *slog.Logger
	DB                *gorm.DB
	TechnologyService *technology.Service
}

package controller

import (
	"mkit/example/ginapp/config"
	"mkit/example/ginapp/internal/service/technology"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DependencyContainer struct {
	Cfg               *config.Config
	Logger            *logrus.Logger
	DB                *gorm.DB
	TechnologyService *technology.Service
}

package postgres

import (
	"fmt"
	"log/slog"
	"mkit/pkg/config"
	"mkit/pkg/enum"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

type DB struct {
	db *gorm.DB
}

func New(slogLogger *slog.Logger, cfg *config.App) (*gorm.DB, error) {
	var (
		dbCfg  = cfg.Postgres
		params = dbCfg.ConnectionParams
		dsn    = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			params.Host, params.Port, params.User, params.Password, params.DBName, params.SSLMode, cfg.Timezone,
		)
		logLevel = logger.Warn
	)

	if cfg.Environment == enum.EnvironmentProduction {
		logLevel = logger.Error
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &GormSlogLogger{
			LogLevel: logLevel,
			Logger:   slogLogger,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("cannot open postgresql connection: %w", err)
	}

	if cfg.Tracing.Enabled {
		if err := db.Use(tracing.NewPlugin()); err != nil {
			return nil, err
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

type CustomMigrator interface {
	CustomMigrate(db *gorm.DB) error
}

func (d *DB) MigrateTables(models ...any) error {
	err := d.db.Migrator().AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("cant migrate tables: %w", err)
	}

	for _, m := range models {
		if migrator, ok := m.(CustomMigrator); ok {
			if err := migrator.CustomMigrate(d.db); err != nil {
				return fmt.Errorf("custom migration failed for %T: %w", m, err)
			}
		}
	}

	return nil
}

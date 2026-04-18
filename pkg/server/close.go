package server

import (
	"context"
	"errors"
	"fmt"

	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
)

func (s *Server) closeInternalServers() {
	logger := s.Deps.Logger

	if len(s.internalGRPCServers) > 0 {
		for _, is := range s.internalGRPCServers {
			name := is.Name()
			logger.Info(fmt.Sprintf("Graceful shutdown for service %s starts", name))
			if err := is.Close(); err != nil && !errors.Is(err, context.Canceled) {
				logger.Warn(fmt.Sprintf("Error when graceful shutdown for service %s: %s", name, err))
			} else {
				logger.Info(fmt.Sprintf("Service %s is shut down successfully", name))
			}
		}
	}

	if s.internalHTTPServer != nil {
		name := s.internalHTTPServer.Name()
		logger.Info(fmt.Sprintf("Graceful shutdown for http server %s starts", name))
		if err := s.internalHTTPServer.Close(); err != nil && !errors.Is(err, context.Canceled) {
			logger.Warn(fmt.Sprintf("Error when graceful shutdown for http server %s: %s", name, err))
		} else {
			logger.Info(fmt.Sprintf("HTTP server %s is shut down successfully", name))
		}
	}

	for _, is := range s.internalConnectServers {
		name := is.Name()
		logger.Info(fmt.Sprintf("Graceful shutdown for connect service %s starts", name))
		if err := is.Close(); err != nil && !errors.Is(err, context.Canceled) {
			logger.Warn(fmt.Sprintf("Error when graceful shutdown for connect service %s: %s", name, err))
		} else {
			logger.Info(fmt.Sprintf("Connect service %s is shut down successfully", name))
		}
	}
}

func (s *Server) Close(ctx context.Context) {
	var (
		d      = s.Deps
		logger = d.Logger
	)

	if healthServer := d.HealthServer; healthServer != nil {
		d.HealthServer.SetServingStatus("", healthv1.HealthCheckResponse_NOT_SERVING)
		healthServer.Shutdown()
		logger.Info("Health server stopped")
	}

	s.closeInternalServers()

	if trace := d.Tracing; trace != nil {
		if err := trace.Stop(ctx); err != nil {
			logger.Warn("Failed to stop tracing", "error", err)
		} else {
			logger.Info("Tracing stopped")
		}
	}

	if gRPCServer := d.GRPCServer; gRPCServer != nil {
		gRPCServer.GracefulStop()
		logger.Info("Grpc server stopped")
	}

	if cronManager := d.CronManager; cronManager != nil {
		if err := cronManager.Stop(); err != nil && !errors.Is(err, context.Canceled) {
			logger.Warn("Failed to stop cron manager", "error", err)
		} else {
			logger.Info("Cron manager stopped")
		}
	}

	if rabbitMQ := d.RabbitMQ; rabbitMQ != nil {
		if err := rabbitMQ.Close(); err != nil {
			logger.Warn("Error when graceful shutdown for RabbitMQ", "error", err)
		} else {
			logger.Info("RabbitMQ shut down")
		}
	}

	if redisClient := d.Redis; redisClient != nil {
		if err := redisClient.Close(); err != nil {
			logger.Warn("Error when graceful shutdown for Redis", "error", err)
		} else {
			logger.Info("Redis client shut down")
		}
	}

	if db := d.Postgres; db != nil {
		if sqlDB, err := db.DB(); err == nil {
			if err = sqlDB.Close(); err != nil {
				logger.Warn("Error when graceful shutdown for Postgresql", "error", err)
			} else {
				logger.Info("Postgresql shut down")
			}
		}
	}
}

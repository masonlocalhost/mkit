package cron

import (
	"context"
	"fmt"
	"mkit/pkg/config"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// Service is a manager to initialize cronjob at microservice startup time
// cron jobs and its configs must not be changed at runtime
type Service struct {
	cron          *cron.Cron
	logger        *logrus.Logger
	redSync       *redsync.Redsync
	location      *time.Location
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
	registry      *sync.Map
}

func New(
	rootCtx context.Context, location *time.Location, logger *logrus.Logger, redSync *redsync.Redsync,
) *Service {
	ctx, cancelFunc := context.WithCancel(context.WithoutCancel(rootCtx))

	s := &Service{
		cron:          cron.New(cron.WithLocation(location)),
		logger:        logger,
		redSync:       redSync,
		location:      location,
		ctx:           ctx,
		ctxCancelFunc: cancelFunc,
		registry:      new(sync.Map),
	}

	s.cron.Start()

	return s
}

func (s *Service) Stop() error {
	s.ctxCancelFunc()
	ctx := s.cron.Stop()
	<-ctx.Done()

	return ctx.Err()
}

func (s *Service) ScheduleCron(cfg *config.Cronjob, handler Handler) error {
	if _, ok := s.registry.Load(cfg.ID); ok {
		return fmt.Errorf("cron [%s] is already scheduled", cfg.ID)
	}
	c := NewCronjob(cfg.ID, cfg.TaskTimeout, s.logger, s.redSync, &handler)

	if cfg.Disabled {
		s.logger.Infof("Skip schedule cronjob [%s] because it is disabled.", cfg.ID)

		return nil
	}

	if _, err := s.cron.AddFunc(cfg.Spec, func() {
		c.Run(s.ctx)
	}); err != nil {
		return err
	}
	s.registry.Store(cfg.ID, struct{}{})
	s.logger.Infof("Cron [%s] is scheduled with spec '%s'", cfg.ID, cfg.Spec)

	return nil
}

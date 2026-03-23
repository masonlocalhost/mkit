package cron

import (
	"context"
	"fmt"
	redsyncutil "mkit/pkg/cache/redsync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

const defaultTaskTimeout = time.Minute * 2

type Handler = func(ctx context.Context) error

// Cronjob is a job to run regularly
// ID is unique string represent it, eg: A:CRONJOB:MANAGE_COLLECTION_AUTOSTART, it's also used for redsync mutex lock name
type Cronjob struct {
	ID          string
	logger      *logrus.Logger
	taskTimeout time.Duration
	handler     *Handler
	redSync     *redsync.Redsync
}

func NewCronjob(
	id string, taskTimeout time.Duration, logger *logrus.Logger, redSync *redsync.Redsync, handler *Handler,
) *Cronjob {
	if taskTimeout == 0 {
		taskTimeout = defaultTaskTimeout
	}

	return &Cronjob{
		ID:          id,
		logger:      logger,
		taskTimeout: taskTimeout,
		handler:     handler,
		redSync:     redSync,
	}
}

// Run wraps Handler for tracing, single pod checking, logging
func (c *Cronjob) Run(sctx context.Context) {
	var (
		mutex = c.redSync.NewMutex(
			c.ID, redsync.WithExpiry(15*time.Second), // lock expiration
			redsync.WithTries(5), // retry attempts
			redsync.WithRetryDelay(200*time.Millisecond),
		)
		tr           = otel.Tracer(fmt.Sprintf("mkit.pkg.cron"))
		rctx, cancel = context.WithTimeout(sctx, c.taskTimeout)
		ctx, span    = tr.Start(rctx, fmt.Sprintf("Cron.%s", c.ID))
		logger       = c.logger.WithContext(ctx).WithField("cron ID", c.ID)
	)

	defer cancel()
	defer span.End()

	ok, err := redsyncutil.TryLock(ctx, mutex)
	if err != nil {
		logger.Errorf("Cannot aqquire lock for check cronjob [%s]: %v.", c.ID, err)

		return
	}
	if !ok {
		logger.Infof("Cronjob [%s] is already running, skipped.", c.ID)

		return
	}

	if err := (*c.handler)(ctx); err != nil {
		logger.Errorf("Failed to execute task of cronjob [%s]: %v", c.ID, err)

		return
	}

	logger.Infof("Cronjob [%s] executed successfully", c.ID)
}

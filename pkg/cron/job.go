package cron

import (
	"context"
	"fmt"
	"log/slog"
	redsyncutil "mkit/pkg/cache/redsync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"go.opentelemetry.io/otel"
)

const defaultTaskTimeout = time.Minute * 2

type Handler = func(ctx context.Context) error

// Cronjob runs a handler regularly.
// ID is a unique string (e.g. "A:CRONJOB:MANAGE_COLLECTION_AUTOSTART") also used as the Redsync mutex name.
type Cronjob struct {
	ID          string
	logger      *slog.Logger
	taskTimeout time.Duration
	handler     *Handler
	redSync     *redsync.Redsync
}

func NewCronjob(
	id string, taskTimeout time.Duration, logger *slog.Logger, redSync *redsync.Redsync, handler *Handler,
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

// Run wraps Handler with tracing, distributed locking, and logging.
func (c *Cronjob) Run(sctx context.Context) {
	var (
		mutex = c.redSync.NewMutex(
			c.ID,
			redsync.WithExpiry(15*time.Second),
			redsync.WithTries(5),
			redsync.WithRetryDelay(200*time.Millisecond),
		)
		tr           = otel.Tracer("mkit.pkg.cron")
		rctx, cancel = context.WithTimeout(sctx, c.taskTimeout)
		ctx, span    = tr.Start(rctx, fmt.Sprintf("Cron.%s", c.ID))
		logger       = c.logger.With("cron_id", c.ID)
	)

	defer cancel()
	defer span.End()

	ok, err := redsyncutil.TryLock(ctx, mutex)
	if err != nil {
		logger.ErrorContext(ctx, "Cannot acquire lock for cronjob", "error", err)
		return
	}
	if !ok {
		logger.InfoContext(ctx, "Cronjob already running, skipped")
		return
	}

	if err := (*c.handler)(ctx); err != nil {
		logger.ErrorContext(ctx, "Failed to execute cronjob task", "error", err)
		return
	}

	logger.InfoContext(ctx, "Cronjob executed successfully")
}

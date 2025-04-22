package task

import (
	"context"
	"log/slog"
	"time"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/service"
)

type AccrualCheckerTask struct {
	config  *config.Config
	service *service.BalanceService
	logger  *slog.Logger
}

func NewAccrualCheckerTask(c *config.Config, s *service.BalanceService, l *slog.Logger) *AccrualCheckerTask {
	return &AccrualCheckerTask{config: c, service: s, logger: l}
}

func (t *AccrualCheckerTask) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(3 * time.Second):
			err := t.service.ProcessPendingOrders(ctx)
			if err != nil {
				t.logger.ErrorContext(ctx, "Error processing task", "err", err.Error())
			}
		}
	}
}

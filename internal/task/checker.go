package task

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/service"
)

type AccrualCheckerTask struct {
	config  *config.Config
	service *service.BalanceService
}

func NewAccrualCheckerTask(config *config.Config, service *service.BalanceService) *AccrualCheckerTask {
	return &AccrualCheckerTask{config: config, service: service}
}

func (t *AccrualCheckerTask) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(3 * time.Second):
			err := t.service.ProcessPendingOrders(ctx)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

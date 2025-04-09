package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
)

type AccrualCheckerTask struct {
	config     *config.Config
	repository repository.Repository
}

func NewAccrualCheckerTask(config *config.Config, repository repository.Repository) *AccrualCheckerTask {
	return &AccrualCheckerTask{config: config, repository: repository}
}

func (t *AccrualCheckerTask) CheckOrderStatus(ctx context.Context, number string) (*models.AccrualStatusDTO, error) {

	url := fmt.Sprintf("%s/api/orders/%s", t.config.AccrualSystemAddress, number)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, common.ErrorNotFound
	}

	reply, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var o *models.AccrualStatusDTO
	err = json.Unmarshal(reply, &o)

	if err != nil {
		return nil, err
	}

	return o, nil

}

func (t *AccrualCheckerTask) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(3 * time.Second):
			orders, err := t.repository.GetUnprocessedOrders(ctx)
			if err != nil {
				fmt.Println(0, err)
				continue
			}
			for _, o := range orders {

				// - `NEW` — заказ загружен в систему, но не попал в обработку;
				// - `PROCESSING` — вознаграждение за заказ рассчитывается;
				// - `INVALID` — система расчёта вознаграждений отказала в расчёте;
				// - `PROCESSED` — данные по заказу проверены и информация о расчёте успешно получена.

				// - `REGISTERED` — заказ зарегистрирован, но не начисление не рассчитано;
				// - `INVALID` — заказ не принят к расчёту, и вознаграждение не будет начислено;
				// - `PROCESSING` — расчёт начисления в процессе;
				// - `PROCESSED` — расчёт начисления окончен;

				fmt.Println("checking status", o)
				accrual, err := t.CheckOrderStatus(ctx, o.Number)
				if err != nil {
					if errors.Is(err, common.ErrorNotFound) {
						fmt.Println("NOT FOUND")

					} else {
						fmt.Println(3, err)
					}
					continue
				}

				var newStatus models.OrderStatus
				var accrualAmount float32

				switch accrual.Status {
				case models.AccrualStatusProcessing:
					newStatus = models.OrderStatusProcessing
				case models.AccrualStatusProcessed:
					newStatus = models.OrderStatusProcessed
					accrualAmount = accrual.Accrual
				case models.AccrualStatusInvalid:
					newStatus = models.OrderStatusInvalid
				}

				fmt.Println("updating status", o)
				_, err = t.repository.UpdateOrderAccrualStatus(ctx, o.ID, newStatus, accrualAmount)

				if err != nil {
					fmt.Println(4, err)
				} else {
					if newStatus == models.OrderStatusProcessed || newStatus == models.OrderStatusInvalid {

					}
				}
			}

		}
	}
}

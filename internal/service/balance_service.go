package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/logging"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
)

type BalanceService struct {
	repository repository.Repository
	config     *config.Config
	logger     logging.Logger
}

func NewBalanceService(r repository.Repository, c *config.Config, l logging.Logger) *BalanceService {
	return &BalanceService{repository: r, config: c, logger: l.With("task", "process_pending_orders")}
}

func (s *BalanceService) checkOrderStatusInAccrualSystem(ctx context.Context, number string) (*models.AccrualStatusDTO, error) {

	url := fmt.Sprintf("%s/api/orders/%s", s.config.AccrualSystemAddress, number)

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

func (s *BalanceService) processOrder(ctx context.Context, order models.Order) error {

	logger := s.logger.With("number", order.Number)

	accrual, err := s.checkOrderStatusInAccrualSystem(ctx, order.Number)
	if err != nil {
		return err
	}

	logger.InfoContext(ctx, "Status received", "status", accrual.Status)

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

	logger.InfoContext(ctx, "Udating status", "status", newStatus)
	_, err = s.repository.UpdateOrderAccrualStatus(ctx, order.ID, newStatus, accrualAmount)

	if err != nil {
		return err
	}

	return s.recalculateAccruals(ctx, order.UserID)
}

func (s *BalanceService) recalculateAccruals(ctx context.Context, userID string) error {

	orders, err := s.repository.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return err
	}
	var totalAccrued float32
	for _, o := range orders {
		totalAccrued += o.Accrual
	}

	return s.repository.UpdateUserAccruedTotel(ctx, userID, totalAccrued)

}

func (s *BalanceService) ProcessPendingOrders(ctx context.Context) error {

	orders, err := s.repository.GetUnprocessedOrders(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error selecting orders", "err", err.Error())
		return err
	}

	for _, o := range orders {

		err := s.processOrder(ctx, o)
		if err != nil {
			s.logger.ErrorContext(ctx, "Error processig order", "number", o.Number, "err", err)
		}
	}

	return nil

}

func (s *BalanceService) GetUserBalance(ctx context.Context, userID string) (*models.BalanceDTO, error) {
	user, err := s.repository.FindUserById(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error finding user", "id", userID, "err", err.Error())
		return nil, err
	}

	return &models.BalanceDTO{Current: user.AccruedTotal - user.WithdrawnTotal, Accrual: user.AccruedTotal}, nil
}

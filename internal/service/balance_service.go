package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

	return &models.BalanceDTO{Current: user.AccruedTotal - float32(user.WithdrawnTotal), Accrual: user.AccruedTotal}, nil
}

func (s *BalanceService) recalculateWithdrawals(ctx context.Context, userID string) error {

	withdrawals, err := s.repository.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return err
	}
	var totalWithdrawn int32
	for _, w := range withdrawals {
		totalWithdrawn += w.Amount
	}

	return s.repository.UpdateUserWithdrawnTotel(ctx, userID, totalWithdrawn)

}

func (s *BalanceService) Withdraw(ctx context.Context, userID string, request *models.WithdrawalRequestDTO) error {

	correct, err := common.CheckOrderNumberFormat(request.Order)
	if err != nil || !correct {
		s.logger.ErrorContext(ctx, "Invalid order number", "number", request.Order)
		return common.ErrorInvalidOrderNumberFormat
	}

	// checking the balance
	user, err := s.repository.FindUserById(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error finding user", "id", userID, "err", err.Error())
		return err
	}

	if user.AccruedTotal-float32(user.WithdrawnTotal)-float32(request.Sum) < 0 {
		s.logger.ErrorContext(ctx, "Insufficient balance", "id", userID)
		return common.ErrorInsufficientBalance
	}

	// user has enough points, making withdrawal
	w := &models.Withdrawal{UploadedAt: time.Now().Truncate(time.Second), UserID: userID, Order: request.Order, Amount: request.Sum}

	_, err = s.repository.AddWithdrawal(ctx, w)

	if err != nil {
		s.logger.ErrorContext(ctx, "Error saving withdrawal", "id", userID, "err", err.Error())
		return err
	}

	return s.recalculateWithdrawals(ctx, userID)

}

func (s *BalanceService) GetWithdrawals(ctx context.Context, userID string) ([]*models.WithdrawalDTO, error) {
	withdrawals, err := s.repository.GetWithdrawalsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []*models.WithdrawalDTO

	for _, w := range withdrawals {
		result = append(result, &models.WithdrawalDTO{Order: w.Order, Sum: w.Amount, ProcessedAt: w.UploadedAt})
	}

	return result, nil
}

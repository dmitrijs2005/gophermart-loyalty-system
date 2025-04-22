package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
)

type BalanceService struct {
	BaseService
	repository repository.Repository
	config     *config.Config
	logger     *slog.Logger
}

func NewBalanceService(r repository.Repository, c *config.Config, l *slog.Logger) *BalanceService {
	return &BalanceService{repository: r, config: c, logger: l.With("task", "process_pending_orders")}
}

func (s *BalanceService) checkOrderStatusInAccrualSystem(ctx context.Context, number string) (*models.AccrualStatusDTO, error) {

	url := fmt.Sprintf("%s/api/orders/%s", s.config.AccrualSystemAddress, number)

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

	tx, err := s.repository.UnitOfWork().Begin(ctx)
	if err != nil {
		return err
	}
	defer s.EndTransaction(tx, &err)

	err = s.repository.UpdateOrderAccrualStatus(ctx, order.ID, newStatus, accrualAmount)

	if err != nil {
		return err
	}

	return s.recalculateAccruals(ctx, order.UserID)
}

func (s *BalanceService) recalculateAccruals(ctx context.Context, userID string) error {

	totalAccrued, err := s.repository.GetAccrualsTotalAmountByUserID(ctx, userID)
	if err != nil {
		return err
	}
	s.logger.With("user_id", userID).InfoContext(ctx, "Updating accrued total", "amount", totalAccrued)

	return s.repository.UpdateUserAccruedTotal(ctx, userID, totalAccrued)

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

			if errors.Is(err, common.ErrorNotFound) {
				s.logger.InfoContext(ctx, "Order not registered in accrual system yet", "number", o.Number)
			} else {
				s.logger.ErrorContext(ctx, "Error processig order", "number", o.Number, "err", err)
			}
		}
	}

	return nil

}

func (s *BalanceService) GetUserBalance(ctx context.Context, userID string) (*models.BalanceDTO, error) {
	user, err := s.repository.FindUserByID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error finding user", "id", userID, "err", err.Error())
		return nil, err
	}

	return &models.BalanceDTO{Current: user.AccruedTotal - user.WithdrawnTotal, Withdrawn: user.WithdrawnTotal}, nil
}

func (s *BalanceService) recalculateWithdrawals(ctx context.Context, userID string) error {

	totalWithdrawn, err := s.repository.GetWithdrawalsTotalAmountByUserID(ctx, userID)
	if err != nil {
		return err
	}

	s.logger.With("user_id", userID).InfoContext(ctx, "Updating withdrawn total", "amount", totalWithdrawn)

	return s.repository.UpdateUserWithdrawnTotal(ctx, userID, totalWithdrawn)

}

func (s *BalanceService) Withdraw(ctx context.Context, userID string, request *models.WithdrawalRequestDTO) error {

	tx, err := s.repository.UnitOfWork().Begin(ctx)
	if err != nil {
		return err
	}
	defer s.EndTransaction(tx, &err)

	correct, err := common.CheckOrderNumberFormat(request.Order)
	if err != nil || !correct {
		s.logger.ErrorContext(ctx, "Invalid order number", "number", request.Order)
		return common.ErrorInvalidOrderNumberFormat
	}

	// checking the balance
	user, err := s.repository.FindUserByID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error finding user", "id", userID, "err", err.Error())
		return err
	}

	if user.AccruedTotal-user.WithdrawnTotal-request.Sum < 0 {
		s.logger.ErrorContext(ctx, "Insufficient balance", "id", userID)
		return common.ErrorInsufficientBalance
	}

	// user has enough points, making withdrawal
	w := &models.Withdrawal{UploadedAt: time.Now().Truncate(time.Second), UserID: userID, Order: request.Order, Amount: request.Sum}

	err = s.repository.AddWithdrawal(ctx, w)

	if err != nil {
		s.logger.ErrorContext(ctx, "Error saving withdrawal", "id", userID, "err", err.Error())
		return err
	}

	s.logger.With("user_id", userID).Info("Saved withdrawal", "amount", request.Sum)

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

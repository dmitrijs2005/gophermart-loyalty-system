package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
)

type OrderStatus int

const (
	OrderStatusSubmittedByAnotherUser OrderStatus = iota
	OrderStatusInvalidNumberFormat
	OrderStatusSubmittedByThisUser
	OrderStatusInternalError
	OrderStatusAccepted
)

type OrderService struct {
	repository repository.Repository
}

func NewOrderService(r repository.Repository) *OrderService {
	return &OrderService{repository: r}
}

func newOrder(userID string, number string) (*models.Order, error) {
	if userID == "" {
		return nil, errors.New("empty user id")
	}
	if number == "" {
		return nil, errors.New("empty order number")
	}
	return &models.Order{Number: number, UserID: userID, UploadedAt: time.Now(), Status: models.OrderStatusNew}, nil
}

func (s *OrderService) checkOrderNumberFormat(number string) (bool, error) {
	if number == "" {
		return false, nil
	}

	if !common.CheckForAllDigits(number) {
		return false, nil
	}

	valid, err := common.CheckLuhn(number)

	if err != nil {
		return false, err
	}

	if !valid {
		return false, err
	}

	return true, nil
}

func (s *OrderService) RegisterOrderNumber(ctx context.Context, userID string, number string) OrderStatus {

	// optional check
	valid, err := s.checkOrderNumberFormat(number)
	if err != nil {
		fmt.Println(1, err)
		return OrderStatusInternalError
	}

	if !valid {
		fmt.Println(3, err)
		return OrderStatusInvalidNumberFormat
	}

	o, err := newOrder(userID, number)
	if err != nil {
		fmt.Println(2, err)
		return OrderStatusInternalError
	}

	order, err := s.repository.AddOrder(ctx, o)
	if err != nil {
		if errors.Is(err, common.ErrorAlreadyExists) {
			if order.UserID == userID {
				return OrderStatusSubmittedByThisUser
			} else {
				return OrderStatusSubmittedByAnotherUser
			}
		}
	}

	return OrderStatusAccepted

}

func (s *OrderService) GetOrderList(ctx context.Context, userID string) ([]models.Order, error) {

	orders, err := s.repository.GetOrdersByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

package service

import (
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
)

// type OrderStatus int

// const (
// 	OrderStatusSubmittedByAnotherUser OrderStatus = iota
// 	OrderStatusInvalidNumberFormat
// 	OrderStatusSubmittedByThisUser
// 	OrderStatusInternalError
// 	OrderStatusAccepted
// )

type BalanceService struct {
	repository repository.Repository
}

func NewBalanceService(r repository.Repository) *BalanceService {
	return &BalanceService{repository: r}
}

// func (s *OrderService) GetUserBalance(ctx context.Context, userID string) OrderStatus {

// 	// optional check
// 	valid, err := s.checkOrderNumberFormat(number)
// 	if err != nil {
// 		fmt.Println(1, err)
// 		return OrderStatusInternalError
// 	}

// 	if !valid {
// 		fmt.Println(3, err)
// 		return OrderStatusInvalidNumberFormat
// 	}

// 	o, err := newOrder(userID, number)
// 	if err != nil {
// 		fmt.Println(2, err)
// 		return OrderStatusInternalError
// 	}

// 	order, err := s.repository.AddOrder(ctx, o)
// 	if err != nil {
// 		if errors.Is(err, common.ErrorAlreadyExists) {
// 			if order.UserID == userID {
// 				return OrderStatusSubmittedByThisUser
// 			} else {
// 				return OrderStatusSubmittedByAnotherUser
// 			}
// 		}
// 	}

// 	return OrderStatusAccepted

// }

// func (s *OrderService) GetOrderList(ctx context.Context, userID string) ([]models.Order, error) {

// 	orders, err := s.repository.GetOrderListByUserId(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return orders, nil
// }

// func (s *OrderService) GetOrderList(ctx context.Context, userID string) ([]models.Order, error) {

// 	orders, err := s.repository.GetOrderListByUserId(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return orders, nil
// }

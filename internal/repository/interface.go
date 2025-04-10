package repository

import (
	"context"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
)

type Repository interface {

	// transaction related
	BeginTransaction(ctx context.Context) error
	CommitTransaction(ctx context.Context) error
	RollbackTransaction(ctx context.Context) error

	// user related
	AddUser(ctx context.Context, user *models.User) (models.User, error)
	FindUserByLogin(ctx context.Context, login string) (models.User, error)
	GetOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error)

	// order and balance related
	AddOrder(ctx context.Context, order *models.Order) (models.Order, error)
	GetUnprocessedOrders(ctx context.Context) ([]models.Order, error)
	FindOrderByID(ctx context.Context, id string) (models.Order, error)
	UpdateOrderAccrualStatus(ctx context.Context, id string, status models.OrderStatus, accrual float32) (models.Order, error)
}

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
	FindUserById(ctx context.Context, userID string) (models.User, error)

	// order and balance related
	AddOrder(ctx context.Context, order *models.Order) (models.Order, error)
	GetUnprocessedOrders(ctx context.Context) ([]models.Order, error)
	FindOrderByID(ctx context.Context, id string) (models.Order, error)
	UpdateOrderAccrualStatus(ctx context.Context, id string, status models.OrderStatus, accrual float32) (models.Order, error)
	UpdateUserAccruedTotel(ctx context.Context, userID string, amount float32) error
	UpdateUserWithdrawnTotel(ctx context.Context, userID string, amount int32) error
	GetOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error)
	AddWithdrawal(ctx context.Context, withdrawal *models.Withdrawal) (models.Withdrawal, error)
	GetWithdrawalsByUserID(ctx context.Context, userID string) ([]models.Withdrawal, error)
}

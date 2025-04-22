package repository

import (
	"context"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
)

type DBStorage interface {
	RunMigrations(ctx context.Context) error
}

type Repository interface {

	// transaction related
	UnitOfWork() UnitOfWork

	// user related
	AddUser(ctx context.Context, user *models.User) (models.User, error)
	FindUserByLogin(ctx context.Context, login string) (models.User, error)
	FindUserByID(ctx context.Context, userID string) (models.User, error)

	// order and balance related
	AddOrder(ctx context.Context, order *models.Order) (models.Order, error)
	FindOrderByNumber(ctx context.Context, number string) (models.Order, error)

	GetUnprocessedOrders(ctx context.Context) ([]models.Order, error)
	UpdateOrderAccrualStatus(ctx context.Context, id string, status models.OrderStatus, accrual float32) error
	UpdateUserAccruedTotal(ctx context.Context, userID string, amount float32) error
	UpdateUserWithdrawnTotal(ctx context.Context, userID string, amount float32) error
	GetOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error)
	AddWithdrawal(ctx context.Context, withdrawal *models.Withdrawal) error
	GetWithdrawalsByUserID(ctx context.Context, userID string) ([]models.Withdrawal, error)
	GetWithdrawalsTotalAmountByUserID(ctx context.Context, userID string) (float32, error)
	GetAccrualsTotalAmountByUserID(ctx context.Context, userID string) (float32, error)
}

type UnitOfWorkTx interface {
	Commit() error
	Rollback() error
}

type UnitOfWork interface {
	Begin(ctx context.Context) (UnitOfWorkTx, error)
}

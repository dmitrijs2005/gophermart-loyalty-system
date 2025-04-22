package repository

import (
	"context"
	"sort"
	"sync"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/google/uuid"
)

type InMemoryRepository struct {
	mu sync.RWMutex

	inTransaction bool

	users       map[string]models.User
	orders      map[string]models.Order
	withdrawals map[string]models.Withdrawal

	userSnapshot       map[string]models.User
	orderSnapshot      map[string]models.Order
	withdrawalSnapshot map[string]models.Withdrawal
}

func NewInMemoryRepository() (*InMemoryRepository, error) {
	return &InMemoryRepository{
		users:       map[string]models.User{},
		orders:      map[string]models.Order{},
		withdrawals: map[string]models.Withdrawal{},
	}, nil
}

func (r *InMemoryRepository) UnitOfWork() UnitOfWork {
	return &InMemoryUnitOfWork{repository: r}
}

func (r *InMemoryRepository) BeginTransaction() error {

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.inTransaction {
		return common.ErrorAlreadyInTranscation
	}

	// creating copies
	r.userSnapshot = r.users
	r.orderSnapshot = r.orders
	r.withdrawalSnapshot = r.withdrawals

	r.inTransaction = true

	return nil
}

func (r *InMemoryRepository) Commit() error {

	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.inTransaction {
		return common.ErrorNotInTranscation
	}

	r.inTransaction = false
	return nil
}

func (r *InMemoryRepository) Rollback() error {

	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.inTransaction {
		return common.ErrorNotInTranscation
	}

	r.users = r.userSnapshot
	r.orders = r.orderSnapshot
	r.withdrawals = r.withdrawalSnapshot

	r.inTransaction = false
	return nil
}

func (r *InMemoryRepository) findUserIDByLogin(_ context.Context, login string) string {

	users := common.FilterMap[models.User](r.users, func(x models.User) bool {
		return x.Login == login
	})

	if len(users) == 0 {
		return ""
	}

	return users[0].ID
}

func (r *InMemoryRepository) FindUserByLogin(ctx context.Context, login string) (models.User, error) {
	id := r.findUserIDByLogin(ctx, login)
	if id == "" {
		return models.User{}, common.ErrorNotFound
	}
	return r.users[id], nil
}

func (r *InMemoryRepository) AddUser(ctx context.Context, user *models.User) (models.User, error) {
	id := r.findUserIDByLogin(ctx, user.Login)
	if id != "" {
		return r.users[id], common.ErrorLoginAlreadyExists
	}

	id, err := r.newUUID()
	if err != nil {
		return models.User{}, err
	}

	user.ID = id
	r.users[user.ID] = *user

	return *user, nil
}

func (r *InMemoryRepository) newUUID() (string, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (r *InMemoryRepository) FindOrderByNumber(_ context.Context, number string) (models.Order, error) {

	orders := common.FilterMap[models.Order](r.orders, func(x models.Order) bool {
		return x.Number == number
	})

	if len(orders) == 0 {
		return models.Order{}, common.ErrorNotFound
	}

	return orders[0], nil
}

func (r *InMemoryRepository) AddOrder(ctx context.Context, order *models.Order) (models.Order, error) {

	id, err := r.newUUID()
	if err != nil {
		return models.Order{}, err
	}

	order.ID = id
	r.orders[id] = *order

	return *order, nil
}

func (r *InMemoryRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error) {

	orders := common.FilterMap[models.Order](r.orders, func(x models.Order) bool {
		return x.UserID == userID
	})

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAt.After(orders[j].UploadedAt)
	})

	return orders, nil
}

func (r *InMemoryRepository) GetUnprocessedOrders(ctx context.Context) ([]models.Order, error) {

	orders := common.FilterMap[models.Order](r.orders, func(x models.Order) bool {
		return x.Status == models.OrderStatusNew || x.Status == models.OrderStatusProcessing
	})

	return orders, nil
}

func (r *InMemoryRepository) UpdateOrderAccrualStatus(ctx context.Context, orderID string,
	status models.OrderStatus, accrual float32) error {

	o, exist := r.orders[orderID]

	if !exist {
		return common.ErrorNotFound
	}

	o.Status = status
	o.Accrual = accrual

	r.orders[orderID] = o

	return nil

}

func (r *InMemoryRepository) UpdateUserAccruedTotal(ctx context.Context, userID string, amount float32) error {

	user, exist := r.users[userID]

	if !exist {
		return common.ErrorNotFound
	}

	user.AccruedTotal = amount

	r.users[userID] = user

	return nil

}

func (r *InMemoryRepository) GetWithdrawalsTotalAmountByUserID(ctx context.Context, userID string) (float32, error) {
	withdrawals := common.FilterMap[models.Withdrawal](r.withdrawals, func(x models.Withdrawal) bool {
		return x.UserID == userID
	})

	var res float32
	for _, w := range withdrawals {
		res += w.Amount
	}

	return res, nil

}

func (r *InMemoryRepository) UpdateUserWithdrawnTotal(ctx context.Context, userID string, amount float32) error {

	user, exist := r.users[userID]

	if !exist {
		return common.ErrorNotFound
	}

	user.WithdrawnTotal = amount

	r.users[userID] = user

	return nil

}

func (r *InMemoryRepository) FindUserByID(ctx context.Context, userID string) (models.User, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return models.User{}, common.ErrorNotFound
	} else {
		return user, nil
	}

}

func (r *InMemoryRepository) AddWithdrawal(ctx context.Context, item *models.Withdrawal) error {
	id, err := r.newUUID()
	if err != nil {
		return err
	}

	item.ID = id
	r.withdrawals[item.ID] = *item

	return nil
}

func (r *InMemoryRepository) GetWithdrawalsByUserID(ctx context.Context, userID string) ([]models.Withdrawal, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	withdrawals := common.FilterMap[models.Withdrawal](r.withdrawals, func(x models.Withdrawal) bool {
		return x.UserID == userID
	})

	sort.Slice(withdrawals, func(i, j int) bool {
		return withdrawals[i].UploadedAt.After(withdrawals[j].UploadedAt)
	})

	return withdrawals, nil
}

func (r *InMemoryRepository) GetAccrualsTotalAmountByUserID(ctx context.Context, userID string) (float32, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	orders := common.FilterMap[models.Order](r.orders, func(x models.Order) bool {
		return x.UserID == userID && x.Status == models.OrderStatusProcessed
	})

	var res float32
	for _, w := range orders {
		res += w.Accrual
	}

	return res, nil

}

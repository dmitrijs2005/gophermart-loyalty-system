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
	mu                sync.RWMutex
	users             map[string]models.User
	userLookupByLogin map[string]string
	orders            map[string]models.Order
	withdrawals       map[string]models.Withdrawal
}

func NewInMemoryRepository() (*InMemoryRepository, error) {
	return &InMemoryRepository{
		users:             map[string]models.User{},
		userLookupByLogin: map[string]string{},
		orders:            map[string]models.Order{},
		withdrawals:       map[string]models.Withdrawal{},
	}, nil
}

func (r *InMemoryRepository) BeginTransaction(ctx context.Context) error {
	r.mu.Lock()
	//fmt.Println("BEGIN TRANSACTIOB")
	return nil
}

func (r *InMemoryRepository) CommitTransaction(ctx context.Context) error {
	r.mu.Unlock()
	//fmt.Println("COMMIT TRANSACTIOB")
	return nil
}

func (r *InMemoryRepository) RollbackTransaction(ctx context.Context) error {
	r.mu.Unlock()
	//fmt.Println("ROLLBACK TRANSACTIOB")
	return nil
}

func (r *InMemoryRepository) findUserIDByLogin(_ context.Context, login string) string {
	id, exists := r.userLookupByLogin[login]
	if !exists {
		return ""
	}
	return id
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
	r.userLookupByLogin[user.Login] = user.ID

	return *user, nil
}

func (r *InMemoryRepository) FindOrderByID(ctx context.Context, id string) (models.Order, error) {
	o, exists := r.orders[id]
	if !exists {
		return models.Order{}, common.ErrorOrderDoesNotExist
	}
	return o, nil
}

func (r *InMemoryRepository) newUUID() (string, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (r *InMemoryRepository) findOrderByNumber(_ context.Context, number string) string {

	orders := r.filterOrders(func(o models.Order) bool {
		return o.Number == number
	})

	if len(orders) == 0 {
		return ""
	}

	return orders[0].ID
}

func (r *InMemoryRepository) AddOrder(ctx context.Context, order *models.Order) (models.Order, error) {

	id := r.findOrderByNumber(ctx, order.Number)
	if id != "" {
		return r.orders[id], common.ErrorAlreadyExists
	}

	existingOrder, exists := r.orders[order.Number]
	if exists {
		return existingOrder, common.ErrorOrderAlreadyExists
	}

	id, err := r.newUUID()
	if err != nil {
		return models.Order{}, err
	}

	order.ID = id
	r.orders[id] = *order

	return *order, nil
}

func (r *InMemoryRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	orders := r.filterOrders(func(o models.Order) bool {
		return o.UserID == userID
	})

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAt.After(orders[j].UploadedAt)
	})

	return orders, nil
}

func (r *InMemoryRepository) filterOrders(predicate func(models.Order) bool) []models.Order {
	var result []models.Order
	for _, order := range r.orders {
		if predicate(order) {
			result = append(result, order)
		}
	}
	return result
}

func (r *InMemoryRepository) GetUnprocessedOrders(ctx context.Context) ([]models.Order, error) {

	orders := r.filterOrders(func(o models.Order) bool {
		return o.Status == models.OrderStatusNew || o.Status == models.OrderStatusProcessing
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

func (r *InMemoryRepository) UpdateUserAccruedTotel(ctx context.Context, userID string, amount float32) error {

	user, exist := r.users[userID]

	if !exist {
		return common.ErrorNotFound
	}

	user.AccruedTotal = amount

	r.users[userID] = user

	return nil

}

func (r *InMemoryRepository) UpdateUserWithdrawnTotel(ctx context.Context, userID string, amount int32) error {

	user, exist := r.users[userID]

	if !exist {
		return common.ErrorNotFound
	}

	user.WithdrawnTotal = amount

	r.users[userID] = user

	return nil

}

func (r *InMemoryRepository) FindUserByID(ctx context.Context, userID string) (models.User, error) {
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

func (r *InMemoryRepository) filterWithdrawals(predicate func(models.Withdrawal) bool) []models.Withdrawal {
	var result []models.Withdrawal
	for _, withdrawal := range r.withdrawals {
		if predicate(withdrawal) {
			result = append(result, withdrawal)
		}
	}
	return result
}

func (r *InMemoryRepository) GetWithdrawalsByUserID(ctx context.Context, userID string) ([]models.Withdrawal, error) {

	withdrawals := r.filterWithdrawals(func(o models.Withdrawal) bool {
		return o.UserID == userID
	})

	sort.Slice(withdrawals, func(i, j int) bool {
		return withdrawals[i].UploadedAt.After(withdrawals[j].UploadedAt)
	})

	return withdrawals, nil
}

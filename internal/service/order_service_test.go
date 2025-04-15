package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/logging"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func TestOrderService_RegisterOrderNumber(t *testing.T) {

	ctx := context.Background()
	repo, err := repository.NewInMemoryRepository()
	require.NoError(t, err)

	user, err := repo.AddUser(ctx, &models.User{Login: "login", Password: "password"})
	require.NoError(t, err)

	user2, err := repo.AddUser(ctx, &models.User{Login: "login2", Password: "password2"})
	require.NoError(t, err)

	config := &config.Config{SecretKey: "secretkey", TokenValidityDuration: 1 * time.Minute}
	logger := logging.NewLogger()

	s := &OrderService{
		repository: repo,
		config:     config,
		logger:     logger,
	}

	type args struct {
		userID string
		number string
	}
	tests := []struct {
		name string
		args args
		want OrderStatus
	}{
		{"OK", args{user.ID, "4561261212345467"}, OrderStatusAccepted},
		{"Wrong number", args{user.ID, "123"}, OrderStatusInvalidNumberFormat},
		{"Same user", args{user.ID, "4561261212345467"}, OrderStatusSubmittedByThisUser},
		{"Another user", args{user2.ID, "4561261212345467"}, OrderStatusSubmittedByAnotherUser},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := s.RegisterOrderNumber(ctx, tt.args.userID, tt.args.number); got != tt.want {
				t.Errorf("OrderService.RegisterOrderNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_GetOrderList(t *testing.T) {

	ctx := context.Background()
	repo, err := repository.NewInMemoryRepository()
	require.NoError(t, err)

	config := &config.Config{SecretKey: "secretkey", TokenValidityDuration: 1 * time.Minute}
	logger := logging.NewLogger()

	user, err := repo.AddUser(ctx, &models.User{Login: "login", Password: "password"})
	require.NoError(t, err)

	user2, err := repo.AddUser(ctx, &models.User{Login: "login2", Password: "password2"})
	require.NoError(t, err)

	order1, err := repo.AddOrder(ctx, &models.Order{UserID: user.ID, Number: "123"})
	require.NoError(t, err)

	order2, err := repo.AddOrder(ctx, &models.Order{UserID: user.ID, Number: "234"})
	require.NoError(t, err)

	s := &OrderService{
		repository: repo,
		config:     config,
		logger:     logger,
	}

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Order
		wantErr bool
	}{
		{"User1", args{ctx, user.ID}, []models.Order{order1, order2}, false},
		{"User2", args{ctx, user2.ID}, []models.Order{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetOrderList(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.GetOrderList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, len(got), len(tt.want))
			if len(tt.want) > 0 {
				if !reflect.DeepEqual(got, tt.want) {
					fmt.Println(got, tt.want)
					t.Errorf("OrderService.GetOrderList() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

package service

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/logging"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestBalanceService(t *testing.T) {

	ctx := context.Background()

	repo, err := repository.NewInMemoryRepository()
	require.NoError(t, err)

	config := &config.Config{SecretKey: "secretkey", TokenValidityDuration: 1 * time.Minute}
	logger := logging.NewLogger()

	s := &BalanceService{
		repository: repo,
		config:     config,
		logger:     logger,
	}

	// setting up data
	user, err := repo.AddUser(ctx, &models.User{Login: "login", Password: "password", AccruedTotal: 5, WithdrawnTotal: 3})
	require.NoError(t, err)
	require.NotZero(t, user.ID)

	user2, err := repo.AddUser(ctx, &models.User{Login: "login2", Password: "password2", AccruedTotal: 8, WithdrawnTotal: 2})
	require.NoError(t, err)
	require.NotZero(t, user.ID)

	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.BalanceDTO
		wantErr bool
	}{
		{"User1", args{user.ID}, &models.BalanceDTO{Withdrawn: user.WithdrawnTotal, Current: user.AccruedTotal - user.WithdrawnTotal}, false},
		{"User2", args{user2.ID}, &models.BalanceDTO{Withdrawn: user2.WithdrawnTotal, Current: user2.AccruedTotal - user2.WithdrawnTotal}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetUserBalance(ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceService.GetUserBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BalanceService.GetUserBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBalanceService_Withdraw(t *testing.T) {
	ctx := context.Background()

	repo, err := repository.NewInMemoryRepository()
	require.NoError(t, err)

	config := &config.Config{SecretKey: "secretkey", TokenValidityDuration: 1 * time.Minute}
	logger := logging.NewLogger()

	s := &BalanceService{
		repository: repo,
		config:     config,
		logger:     logger,
	}

	// setting up data
	user, err := repo.AddUser(ctx, &models.User{Login: "login", Password: "password", AccruedTotal: 5, WithdrawnTotal: 3})
	require.NoError(t, err)
	require.NotZero(t, user.ID)

	type args struct {
		userID  string
		request *models.WithdrawalRequestDTO
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"OK", args{user.ID, &models.WithdrawalRequestDTO{Order: "4561261212345467", Sum: float32(1)}}, false},
		{"Wrong format", args{user.ID, &models.WithdrawalRequestDTO{Order: "123", Sum: float32(1)}}, true},
		{"Insufficent balance", args{user.ID, &models.WithdrawalRequestDTO{Order: "4561261212345467", Sum: float32(10)}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Withdraw(ctx, tt.args.userID, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("BalanceService.Withdraw() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBalanceService_GetWithdrawals(t *testing.T) {
	ctx := context.Background()

	repo, err := repository.NewInMemoryRepository()
	require.NoError(t, err)

	config := &config.Config{SecretKey: "secretkey", TokenValidityDuration: 1 * time.Minute}
	logger := logging.NewLogger()

	s := &BalanceService{
		repository: repo,
		config:     config,
		logger:     logger,
	}

	// setting up data
	user, err := repo.AddUser(ctx, &models.User{Login: "login", Password: "password"})
	require.NoError(t, err)
	require.NotZero(t, user.ID)

	err = repo.AddWithdrawal(ctx, &models.Withdrawal{UserID: user.ID, Amount: 1, Order: "123"})
	require.NoError(t, err)

	err = repo.AddWithdrawal(ctx, &models.Withdrawal{UserID: user.ID, Amount: 2, Order: "345"})
	require.NoError(t, err)

	x1 := models.WithdrawalDTO{Order: "123", Sum: 1}
	x2 := models.WithdrawalDTO{Order: "345", Sum: 2}

	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*models.WithdrawalDTO
		wantErr bool
	}{
		{"OK", args{user.ID}, []*models.WithdrawalDTO{&x1, &x2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := s.GetWithdrawals(ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceService.GetWithdrawals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.True(t, cmp.Equal(got, tt.want, cmp.AllowUnexported(models.WithdrawalDTO{})))
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("BalanceService.GetWithdrawals() = %v, want %v", got, tt.want)
			// }
		})
	}
}

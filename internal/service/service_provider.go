package service

import (
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/logging"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
)

type ServiceProvider struct {
	AuthService    *AuthService
	OrderService   *OrderService
	BalanceService *BalanceService
}

func NewServiceProvider(repository repository.Repository, config *config.Config, logger logging.Logger) *ServiceProvider {

	authService := NewAuthService(repository, config, logger)
	orderService := NewOrderService(repository, config, logger)
	balanceService := NewBalanceService(repository, config, logger)

	return &ServiceProvider{AuthService: authService, OrderService: orderService, BalanceService: balanceService}
}

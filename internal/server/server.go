package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/logging"
	m "github.com/dmitrijs2005/gophermart-loyalty-system/internal/server/middleware"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HTTPServer struct {
	config          *config.Config
	logger          logging.Logger
	serviceProvider *service.ServiceProvider
}

func NewHTTPServer(c *config.Config, sp *service.ServiceProvider, logger logging.Logger) (*HTTPServer, error) {
	return &HTTPServer{config: c, serviceProvider: sp, logger: logger}, nil
}

func (s *HTTPServer) RegisterAuthRoutes(r chi.Router) {

	service := s.serviceProvider.AuthService
	h := NewAuthHandler(service)

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
}

func (s *HTTPServer) RegisterOrderRoutes(r chi.Router) {

	service := s.serviceProvider.OrderService
	h := NewOrderHandler(service)

	r.Group(func(r chi.Router) {
		r.Use(m.NewAuthMiddleware(s.config.SecretKey))
		r.Post("/orders", h.RegisterOrder)
		r.Get("/orders", h.GetUserOrderList)
	})

}

func (s *HTTPServer) RegisterBalanceRoutes(r chi.Router) {

	service := s.serviceProvider.BalanceService
	_ = NewBalanceHandler(service)

	r.Group(func(r chi.Router) {
		r.Use(m.NewAuthMiddleware(s.config.SecretKey))
		// r.Post("/orders", h.RegisterOrder)
		// r.Get("/orders", h.GetUserOrderList)
	})

}

func (s *HTTPServer) RegisterRoutes() http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// * `POST /api/user/register` — регистрация пользователя;
	// * `POST /api/user/login` — аутентификация пользователя;
	// * `POST /api/user/orders` — загрузка пользователем номера заказа для расчёта;
	// * `GET /api/user/orders` — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
	// * `GET /api/user/balance` — получение текущего баланса счёта баллов лояльности пользователя;
	// * `POST /api/user/balance/withdraw` — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
	// * `GET /api/user/withdrawals` — получение информации о выводе средств с накопительного счёта пользователем.

	r.Route("/api/user", func(r chi.Router) {
		s.RegisterAuthRoutes(r)
		s.RegisterOrderRoutes(r)
		s.RegisterBalanceRoutes(r)
	})

	return r
}

func (s *HTTPServer) Run(ctx context.Context) error {

	mux := s.RegisterRoutes()

	server := &http.Server{
		Addr:    s.config.RunAddress,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			s.logger.ErrorContext(ctx, "Error", "err", err.Error())
		}
	}()

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

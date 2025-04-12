package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	m "github.com/dmitrijs2005/gophermart-loyalty-system/internal/server/middleware"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HTTPServer struct {
	config          *config.Config
	logger          *slog.Logger
	serviceProvider *service.ServiceProvider
}

func NewHTTPServer(c *config.Config, sp *service.ServiceProvider, logger *slog.Logger) (*HTTPServer, error) {
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
	h := NewOrderHandler(service, s.logger)

	r.Group(func(r chi.Router) {
		r.Use(m.NewAuthMiddleware(s.config.SecretKey))
		r.Post("/orders", h.RegisterOrder)
		r.Get("/orders", h.GetUserOrderList)
	})

}

func (s *HTTPServer) RegisterBalanceRoutes(r chi.Router) {

	service := s.serviceProvider.BalanceService
	h := NewBalanceHandler(service)

	r.Group(func(r chi.Router) {
		r.Use(m.NewAuthMiddleware(s.config.SecretKey))
		r.Get("/balance", h.UserBalance)
		r.Post("/balance/withdraw", h.Withdraw)
		r.Get("/withdrawals", h.Withdrawals)
	})

}

func (s *HTTPServer) RegisterRoutes() http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

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

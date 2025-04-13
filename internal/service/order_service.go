package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
)

type OrderStatus int

const (
	OrderStatusSubmittedByAnotherUser OrderStatus = iota
	OrderStatusInvalidNumberFormat
	OrderStatusSubmittedByThisUser
	OrderStatusInternalError
	OrderStatusAccepted
)

type OrderService struct {
	BaseService
	repository repository.Repository
	config     *config.Config
	logger     *slog.Logger
}

func NewOrderService(r repository.Repository, c *config.Config, l *slog.Logger) *OrderService {
	return &OrderService{repository: r, config: c, logger: l}
}

func newOrder(userID string, number string) (*models.Order, error) {
	if userID == "" {
		return nil, errors.New("empty user id")
	}
	if number == "" {
		return nil, errors.New("empty order number")
	}
	return &models.Order{Number: number, UserID: userID, UploadedAt: time.Now().Truncate(time.Second), Status: models.OrderStatusNew}, nil
}

// #### **Загрузка номера заказа**
// Хендлер: `POST /api/user/orders`.
// Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.
// Номер заказа может быть проверен на корректность ввода с помощью [алгоритма Луна](https://ru.wikipedia.org/wiki/Алгоритм_Луна){target="_blank"}.
// Формат запроса:
// ```
// POST /api/user/orders HTTP/1.1
// Content-Type: text/plain
// ...
// 12345678903
// ```
// Возможные коды ответа:
// - `200` — номер заказа уже был загружен этим пользователем;
// - `202` — новый номер заказа принят в обработку;
// - `400` — неверный формат запроса;
// - `401` — пользователь не аутентифицирован;
// - `409` — номер заказа уже был загружен другим пользователем;
// - `422` — неверный формат номера заказа;
// - `500` — внутренняя ошибка сервера.

func (s *OrderService) RegisterOrderNumber(ctx context.Context, userID string, number string) OrderStatus {

	tx, err := s.repository.UnitOfWork().Begin(ctx)
	if err != nil {
		return OrderStatusInternalError
	}

	defer s.EndTransaction(tx, &err)

	valid, err := common.CheckOrderNumberFormat(number)
	if err != nil {
		return OrderStatusInternalError
	}
	if !valid {
		return OrderStatusInvalidNumberFormat
	}

	existingOrder, err := s.repository.FindOrderByNumber(ctx, number)
	if err != nil {
		if !errors.Is(err, common.ErrorNotFound) {
			return OrderStatusInternalError
		}
	} else {
		if existingOrder.UserID == userID {
			return OrderStatusSubmittedByThisUser
		} else {
			return OrderStatusSubmittedByAnotherUser
		}
	}

	o, err := newOrder(userID, number)
	if err != nil {
		return OrderStatusInternalError
	}

	_, err = s.repository.AddOrder(ctx, o)
	if err != nil {
		return OrderStatusInternalError
	}

	return OrderStatusAccepted

}

func (s *OrderService) GetOrderList(ctx context.Context, userID string) ([]models.Order, error) {

	orders, err := s.repository.GetOrdersByUserID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return nil, err
	}

	return orders, nil

}

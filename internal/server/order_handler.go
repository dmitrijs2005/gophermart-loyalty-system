package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/server/middleware"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
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

func (h *OrderHandler) RegisterOrder(w http.ResponseWriter, r *http.Request) {

	// reading the request body which should contain the order number
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// trying to get userid from context
	a := ctx.Value(middleware.UserIDKey)
	userID, ok := a.(string)
	if !ok {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return

	}

	orderNumber := string(body)

	status := h.service.RegisterOrderNumber(ctx, userID, orderNumber)

	switch status {
	case service.OrderStatusSubmittedByAnotherUser:
		http.Error(w, "order submitted by another user", http.StatusConflict)
		return
	case service.OrderStatusSubmittedByThisUser:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("order submitted by this user"))
		return
	case service.OrderStatusAccepted:
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte{})
		return
	case service.OrderStatusInvalidNumberFormat:
		http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
		return
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

}

// #### **Получение списка загруженных номеров заказов**

// Хендлер: `GET /api/user/orders`.

// Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых новых к самым старым. Формат даты — RFC3339.

// Доступные статусы обработки расчётов:

// - `NEW` — заказ загружен в систему, но не попал в обработку;
// - `PROCESSING` — вознаграждение за заказ рассчитывается;
// - `INVALID` — система расчёта вознаграждений отказала в расчёте;
// - `PROCESSED` — данные по заказу проверены и информация о расчёте успешно получена.

// Формат запроса:

// ```
// GET /api/user/orders HTTP/1.1
// Content-Length: 0
// ```

// Возможные коды ответа:

// - `200` — успешная обработка запроса.

//   Формат ответа:

//     ```
//     200 OK HTTP/1.1
//     Content-Type: application/json
//     ...

//     [
//     	{
//             "number": "9278923470",
//             "status": "PROCESSED",
//             "accrual": 500,
//             "uploaded_at": "2020-12-10T15:15:45+03:00"
//         },
//         {
//             "number": "12345678903",
//             "status": "PROCESSING",
//             "uploaded_at": "2020-12-10T15:12:01+03:00"
//         },
//         {
//             "number": "346436439",
//             "status": "INVALID",
//             "uploaded_at": "2020-12-09T16:09:53+03:00"
//         }
//     ]
//     ```

// - `204` — нет данных для ответа.
// - `401` — пользователь не авторизован.
// - `500` — внутренняя ошибка сервера.

func (h *OrderHandler) GetUserOrderList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// trying to get userid from context
	a := ctx.Value(middleware.UserIDKey)
	userID, ok := a.(string)
	if !ok {
		http.Error(w, "sas", http.StatusInternalServerError)
		return

	}

	orders, err := h.service.GetOrderList(ctx, userID)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	var reply []models.OrderDTO
	for _, r := range orders {
		reply = append(reply, models.OrderDTO{
			Number:     r.Number,
			Status:     r.Status,
			UploadedAt: r.UploadedAt,
			Accrual:    r.Accrual,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reply); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

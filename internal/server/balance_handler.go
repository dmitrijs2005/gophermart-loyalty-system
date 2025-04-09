package server

import (
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/service"
)

type BalanceHandler struct {
	service *service.BalanceService
}

func NewBalanceHandler(s *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{service: s}
}

// func (h *OrderHandler) RegisterOrder(w http.ResponseWriter, r *http.Request) {

// 	// reading the request body which should contain the order number
// 	body, err := io.ReadAll(r.Body)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	ctx := r.Context()

// 	// trying to get userid from context
// 	a := ctx.Value("UserID")
// 	userID, ok := a.(string)
// 	if !ok {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return

// 	}

// 	orderNumber := string(body)

// 	status := h.service.RegisterOrderNumber(ctx, userID, orderNumber)

// 	switch status {
// 	case service.OrderStatusSubmittedByAnotherUser:
// 		http.Error(w, "order submitted by another user", http.StatusConflict)
// 		return
// 	case service.OrderStatusSubmittedByThisUser:
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("order submitted by this user"))
// 		return
// 	case service.OrderStatusAccepted:
// 		w.WriteHeader(http.StatusAccepted)
// 		w.Write([]byte{})
// 		return
// 	case service.OrderStatusInvalidNumberFormat:
// 		http.Error(w, "invalid order number format", http.StatusUnprocessableEntity)
// 		return
// 	default:
// 		http.Error(w, "internal error", http.StatusInternalServerError)
// 		return
// 	}

// }

// // #### **Получение списка загруженных номеров заказов**

// // Хендлер: `GET /api/user/orders`.

// // Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых новых к самым старым. Формат даты — RFC3339.

// // Доступные статусы обработки расчётов:

// // - `NEW` — заказ загружен в систему, но не попал в обработку;
// // - `PROCESSING` — вознаграждение за заказ рассчитывается;
// // - `INVALID` — система расчёта вознаграждений отказала в расчёте;
// // - `PROCESSED` — данные по заказу проверены и информация о расчёте успешно получена.

// // Формат запроса:

// // ```
// // GET /api/user/orders HTTP/1.1
// // Content-Length: 0
// // ```

// // Возможные коды ответа:

// // - `200` — успешная обработка запроса.

// //   Формат ответа:

// //     ```
// //     200 OK HTTP/1.1
// //     Content-Type: application/json
// //     ...

// //     [
// //     	{
// //             "number": "9278923470",
// //             "status": "PROCESSED",
// //             "accrual": 500,
// //             "uploaded_at": "2020-12-10T15:15:45+03:00"
// //         },
// //         {
// //             "number": "12345678903",
// //             "status": "PROCESSING",
// //             "uploaded_at": "2020-12-10T15:12:01+03:00"
// //         },
// //         {
// //             "number": "346436439",
// //             "status": "INVALID",
// //             "uploaded_at": "2020-12-09T16:09:53+03:00"
// //         }
// //     ]
// //     ```

// // - `204` — нет данных для ответа.
// // - `401` — пользователь не авторизован.
// // - `500` — внутренняя ошибка сервера.

// func (h *OrderHandler) GetUserOrderList(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	// trying to get userid from context
// 	a := ctx.Value("UserID")
// 	userID, ok := a.(string)
// 	if !ok {
// 		http.Error(w, "sas", http.StatusInternalServerError)
// 		return

// 	}

// 	orders, err := h.service.GetOrderList(ctx, userID)
// 	if err != nil {
// 		http.Error(w, "error", http.StatusInternalServerError)
// 		return
// 	}

// 	var reply []models.OrderDTO
// 	for _, r := range orders {
// 		reply = append(reply, models.OrderDTO{
// 			Number:     r.Number,
// 			Status:     r.Status,
// 			UploadedAt: r.UploadedAt,
// 			Accrual:    r.Accrual,
// 		})
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	if err := json.NewEncoder(w).Encode(reply); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// }

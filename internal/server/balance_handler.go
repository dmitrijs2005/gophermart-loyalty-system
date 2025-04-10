package server

import (
	"encoding/json"
	"net/http"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/server/middleware"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/service"
)

type BalanceHandler struct {
	service *service.BalanceService
}

func NewBalanceHandler(s *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{service: s}
}

// #### **Получение текущего баланса пользователя**
// Хендлер: `GET /api/user/balance`.
// Хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов лояльности, а также сумме использованных за весь период регистрации баллов.
// Формат запроса:
// ```
// GET /api/user/balance HTTP/1.1
// Content-Length: 0
// ```
// Возможные коды ответа:
// - `200` — успешная обработка запроса.
//   Формат ответа:
//     ```
//     200 OK HTTP/1.1
//     Content-Type: application/json
//     ...
//     {
//     	"current": 500.5,
//     	"withdrawn": 42
//     }
//     ```
// - `401` — пользователь не авторизован.
// - `500` — внутренняя ошибка сервера.

func (h *BalanceHandler) UserBalance(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// trying to get userid from context
	a := ctx.Value(middleware.UserIDKey)
	userID, ok := a.(string)
	if !ok {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return

	}

	balance, err := h.service.GetUserBalance(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(balance); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

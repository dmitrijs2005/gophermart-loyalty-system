package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/server/middleware"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/service"
	"github.com/go-playground/validator/v10"
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

// #### **Запрос на списание средств**
// Хендлер: `POST /api/user/balance/withdraw`
// Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер нового заказа пользователя в счет оплаты которого списываются баллы.
// Примечание: для успешного списания достаточно успешной регистрации запроса, никаких внешних систем начисления не предусмотрено и не требуется реализовывать.
// Формат запроса:
// ```
// POST /api/user/balance/withdraw HTTP/1.1
// Content-Type: application/json
// {
// 	"order": "2377225624",
//     "sum": 751
// }
// ```
// Здесь `order` — номер заказа, а `sum` — сумма баллов к списанию в счёт оплаты.
// Возможные коды ответа:
// - `200` — успешная обработка запроса;
// - `401` — пользователь не авторизован;
// - `402` — на счету недостаточно средств;
// - `422` — неверный номер заказа;
// - `500` — внутренняя ошибка сервера.

func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {

	var req models.WithdrawalRequestDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	validate := validator.New()
	err = validate.StructCtx(ctx, req)
	if err != nil {
		http.Error(w, common.ErrorValidation.Error(), http.StatusBadRequest)
		return
	}

	// trying to get userid from context
	a := ctx.Value(middleware.UserIDKey)
	userID, ok := a.(string)
	if !ok {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return

	}

	err = h.service.Withdraw(ctx, userID, &req)
	if err != nil {
		if errors.Is(err, common.ErrorInsufficientBalance) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		} else {
			if errors.Is(err, common.ErrorInvalidOrderNumberFormat) {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return

			}
		}
	}

	w.Write([]byte{})

}

// #### **Получение информации о выводе средств**
// Хендлер: `GET /api/user/withdrawals`.
// Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых новых к самым старым. Формат даты — RFC3339.
// Формат запроса:
// ```
// GET /api/user/withdrawals HTTP/1.1
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
//         {
//             "order": "2377225624",
//             "sum": 500,
//             "processed_at": "2020-12-09T16:09:57+03:00"
//         }
//     ]
//     ```
// - `204` - нет ни одного списания.
// - `401` — пользователь не авторизован.
// - `500` — внутренняя ошибка сервера.

func (h *BalanceHandler) Withdrawals(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// trying to get userid from context
	a := ctx.Value(middleware.UserIDKey)
	userID, ok := a.(string)
	if !ok {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return

	}

	result, err := h.service.GetWithdrawals(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(result) == 0 {
		http.Error(w, "no data", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

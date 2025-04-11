package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/service"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// #### **Регистрация пользователя**
// Хендлер: `POST /api/user/register`.
// Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным.
// После успешной регистрации должна происходить автоматическая аутентификация пользователя.
// Формат запроса:
// ```
// POST /api/user/register HTTP/1.1
// Content-Type: application/json
// ...
// {
// 	"login": "<login>",
// 	"password": "<password>"
// }
// ```
// Возможные коды ответа:
// - `200` — пользователь успешно зарегистрирован и аутентифицирован;
// - `400` — неверный формат запроса;
// - `409` — логин уже занят;
// - `500` — внутренняя ошибка сервера.

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req models.RegisterUserDTO

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

	token, err := h.service.Register(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, common.ErrorLoginAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		} else if errors.Is(err, common.ErrorInvalidPasswordFormat) || errors.Is(err, common.ErrorInvalidLoginFormat) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.Write([]byte{})
}

// #### **Аутентификация пользователя**
// Хендлер: `POST /api/user/login`.
// Аутентификация производится по паре логин/пароль.
// Формат запроса:
// ```
// POST /api/user/login HTTP/1.1
// Content-Type: application/json
// ...
// {
// 	"login": "<login>",
// 	"password": "<password>"
// }
// ```
// Возможные коды ответа:
// - `200` — пользователь успешно аутентифицирован;
// - `400` — неверный формат запроса;
// - `401` — неверная пара логин/пароль;
// - `500` — внутренняя ошибка сервера.

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginDTO
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

	token, err := h.service.Login(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, common.ErrorInvalidLoginPassword) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else if errors.Is(err, common.ErrorNotFound) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.Write([]byte{})

}

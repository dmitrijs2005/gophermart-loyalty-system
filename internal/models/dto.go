package models

import (
	"time"
)

type LoginDTO struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterUserDTO struct {
	LoginDTO
}

type OrderDTO struct {
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    float32     `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

type BalanceDTO struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type AccrualStatus string

const (
	AccrualStatusRegistered AccrualStatus = "REGISTERED"
	AccrualStatusInvalid    AccrualStatus = "INVALID"
	AccrualStatusProcessing AccrualStatus = "PROCESSING"
	AccrualStatusProcessed  AccrualStatus = "PROCESSED"
)

type AccrualStatusDTO struct {
	Order   string        `json:"order"`
	Status  AccrualStatus `json:"status"`
	Accrual float32       `json:"accrual"`
}

type WithdrawalRequestDTO struct {
	Order string  `json:"order" validate:"required"`
	Sum   float32 `json:"sum" validate:"required"`
}

type WithdrawalDTO struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

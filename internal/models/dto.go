package models

import "time"

type RegisterUserDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OrderDTO struct {
	ID         string      `json:"id"`
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    float32     `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

type BalanceDTO struct {
	Current float32 `json:"current"`
	Accrual float32 `json:"accrual"`
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
	Order string `json:"order"`
	Sum   int32  `json:"sum"`
}

type WithdrawalDTO struct {
	Order       string    `json:"order"`
	Sum         int32     `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

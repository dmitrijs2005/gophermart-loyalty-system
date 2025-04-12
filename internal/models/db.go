package models

import "time"

type User struct {
	ID             string
	Login          string
	Password       string
	AccruedTotal   float32
	WithdrawnTotal float32
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = `NEW`        //заказ загружен в систему, но не попал в обработку;
	OrderStatusProcessing OrderStatus = `PROCESSING` //вознаграждение за заказ рассчитывается;
	OrderStatusInvalid    OrderStatus = `INVALID`    //система расчёта вознаграждений отказала в расчёте;
	OrderStatusProcessed  OrderStatus = `PROCESSED`  //данные по заказу проверены и информация о расчёте успешно получена.
)

type Order struct {
	ID         string
	Number     string
	UserID     string
	Status     OrderStatus
	Accrual    float32
	UploadedAt time.Time
}

type Withdrawal struct {
	ID         string
	UserID     string
	UploadedAt time.Time
	Order      string
	Amount     float32
}

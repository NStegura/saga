package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderMessageStatus string
type OrderState string

const (
	Created OrderMessageStatus = "CREATED"

	OrderCreated     OrderState = "ORDER_CREATED"
	ReserveCreated   OrderState = "RESERVE_CREATED"
	ReserveFailed    OrderState = "RESERVE_FAILED"
	PaymentCreated   OrderState = "PAYMENT_CREATED"
	PaymentFailed    OrderState = "PAYMENT_FAILED"
	PaymentCompleted OrderState = "PAYMENT_COMPLETED"
)

type OrderMessage struct {
	IKey    uuid.UUID          `json:"idempotent_key"`
	OrderID int64              `json:"order_id"`
	Status  OrderMessageStatus `json:"status"`
}

type Order struct {
	OrderInfo
	OrderProducts []OrderProduct
}

type OrderInfo struct {
	OrderID     int64
	Description string
	State       OrderState
}

type OrderProduct struct {
	ProductID int64
	Count     int64
}

type State struct {
	State     OrderState
	CreatedAt time.Time
}

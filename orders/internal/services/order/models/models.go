package models

import (
	"github.com/google/uuid"
)

type OrderMessageStatus string
type OrderState string

const (
	CREATED OrderMessageStatus = "CREATED"

	ORDER_CREATED     OrderState = "ORDER_CREATED"
	RESERVE_CREATED   OrderState = "RESERVE_CREATED"
	RESERVE_FAILED    OrderState = "RESERVE_FAILED"
	PAYMENT_CREATED   OrderState = "PAYMENT_CREATED"
	PAYMENT_FAILED    OrderState = "PAYMENT_FAILED"
	PAYMENT_COMPLETED OrderState = "PAYMENT_COMPLETED"
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

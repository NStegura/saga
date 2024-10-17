package models

import (
	"github.com/google/uuid"
)

type PaymentMessageStatus string

const (
	CREATED   PaymentMessageStatus = "CREATED"
	FAILED    PaymentMessageStatus = "FAILED"
	COMPLETED PaymentMessageStatus = "COMPLETED"
)

type PaymentMessage struct {
	IKey    uuid.UUID            `json:"idempotent_key"`
	OrderID int64                `json:"order_id"`
	Status  PaymentMessageStatus `json:"status"`
}

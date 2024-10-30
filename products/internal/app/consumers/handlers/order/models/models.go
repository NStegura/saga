package models

import "github.com/google/uuid"

type OrderMessageStatus string

const (
	CREATED OrderMessageStatus = "CREATED"
)

type OrderMessage struct {
	IKey    uuid.UUID          `json:"idempotent_key"`
	OrderID int64              `json:"order_id"`
	Status  OrderMessageStatus `json:"status"`
}

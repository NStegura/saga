package models

import "github.com/google/uuid"

type ProductMessageStatus string

const (
	FAILED  ProductMessageStatus = "FAILED"
	CREATED ProductMessageStatus = "CREATED"
)

type ProductMessage struct {
	IKey    uuid.UUID            `json:"idempotent_key"`
	OrderID int64                `json:"order_id"`
	Status  ProductMessageStatus `json:"status"`
}

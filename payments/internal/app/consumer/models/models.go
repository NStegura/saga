package models

import "github.com/google/uuid"

type InventoryMessageStatus string

const (
	FAILED  InventoryMessageStatus = "FAILED"
	CREATED InventoryMessageStatus = "CREATED"
)

type InventoryMessage struct {
	IKey    uuid.UUID              `json:"idempotent_key"`
	OrderID int64                  `json:"order_id"`
	Status  InventoryMessageStatus `json:"status"`
}

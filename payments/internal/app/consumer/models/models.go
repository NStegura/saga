package models

import "github.com/google/uuid"

type InventoryMessageStatus string

const (
	FAILED    InventoryMessageStatus = "FAILED"
	COMPLETED InventoryMessageStatus = "COMPLETED"
)

type InventoryMessage struct {
	IKey    uuid.UUID              `json:"idempotent_key"`
	OrderID int64                  `json:"order_id"`
	Status  InventoryMessageStatus `json:"status"`
}

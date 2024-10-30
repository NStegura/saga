package models

import (
	"github.com/google/uuid"
)

type ProductMessageStatus string

const (
	CREATED   ProductMessageStatus = "CREATED"
	FAILED    ProductMessageStatus = "FAILED"
	COMPLETED ProductMessageStatus = "COMPLETED"
)

type ProductMessage struct {
	IKey    uuid.UUID            `json:"idempotent_key"`
	OrderID int64                `json:"order_id"`
	Status  ProductMessageStatus `json:"status"`
}

type Product struct {
	ProductID   int64
	Category    string
	Name        string
	Description string
	Count       int64
}

type Reserve struct {
	ProductID int64
	Count     int64
}

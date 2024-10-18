package models

import (
	"database/sql"
)

type Order struct {
	ID          int64          `db:"id"`
	Description sql.NullString `db:"description"`
}

type OrderProduct struct {
	ID        int64 `db:"id"`
	OrderID   int64 `db:"order_id"`
	ProductID int64 `db:"product_id"`
	Count     int64 `db:"count"`
}

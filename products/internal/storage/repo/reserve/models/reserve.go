package models

import (
	"database/sql"
	"time"
)

type Reserve struct {
	ID        int64        `db:"id"`
	ProductID int64        `db:"product_id"`
	OrderID   int64        `db:"order_id"`
	Count     int64        `db:"count"`
	PayStatus sql.NullBool `db:"pay_status"`
	SavedAt   time.Time    `db:"saved_at"`
}

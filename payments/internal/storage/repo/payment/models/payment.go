package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type PaymentStatus int

const (
	CREATED PaymentStatus = iota + 1
	FAILED
	COMPLETED
)

func (ps PaymentStatus) String() string {
	return [...]string{"CREATED", "FAILED", "COMPLETED"}[ps-1]
}

func (ps PaymentStatus) Value() (driver.Value, error) {
	return ps.String(), nil
}

func (ps *PaymentStatus) Scan(value any) error {
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("PaymentStatus should be a string")
	}

	switch strValue {
	case "CREATED":
		*ps = CREATED
	case "FAILED":
		*ps = FAILED
	case "COMPLETED":
		*ps = COMPLETED
	default:
		return fmt.Errorf("invalid PaymentStatus: %s", strValue)
	}
	return nil
}

type Payment struct {
	ID        int64         `db:"id"`
	OrderID   int64         `db:"order_id"`
	Status    PaymentStatus `db:"status"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
}

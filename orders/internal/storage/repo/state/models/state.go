package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type OrderStateStatus int

const (
	OrderCreated OrderStateStatus = iota + 1
	ReserveCreated
	ReserveFailed
	PaymentCreated
	PaymentFailed
	PaymentCompleted
)

func (ps OrderStateStatus) String() string {
	return [...]string{
		"ORDER_CREATED",
		"RESERVE_CREATED", "RESERVE_FAILED",
		"PAYMENT_CREATED", "PAYMENT_FAILED", "PAYMENT_COMPLETED",
	}[ps-1]
}

func (ps OrderStateStatus) Value() (driver.Value, error) {
	return ps.String(), nil
}

func (ps *OrderStateStatus) Scan(value any) error {
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("PaymentStatus should be a string")
	}

	switch strValue {
	case "ORDER_CREATED":
		*ps = OrderCreated
	case "RESERVE_CREATED":
		*ps = ReserveCreated
	case "RESERVE_FAILED":
		*ps = ReserveFailed
	case "PAYMENT_CREATED":
		*ps = PaymentCreated
	case "PAYMENT_FAILED":
		*ps = PaymentFailed
	case "PAYMENT_COMPLETED":
		*ps = PaymentCompleted
	default:
		return fmt.Errorf("invalid OrderStateStatus: %s", strValue)
	}
	return nil
}

type OrderState struct {
	ID        int64            `db:"id"`
	OrderID   int64            `db:"order_id"`
	CreatedAt time.Time        `db:"created_at"`
	State     OrderStateStatus `db:"state"`
}

package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type OutboxStatus int

const (
	WAIT OutboxStatus = iota + 1
	PROCESSED
)

func (os *OutboxStatus) String() string {
	return [...]string{"WAIT", "PROCESSED"}[*os-1]
}

func (os *OutboxStatus) Value() (driver.Value, error) {
	return os.String(), nil
}

func (os *OutboxStatus) Scan(value any) error {
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("OutboxStatus should be a string")
	}

	switch strValue {
	case "WAIT":
		*os = WAIT
	case "PROCESSED":
		*os = PROCESSED
	default:
		return fmt.Errorf("invalid OutboxStatus: %s", strValue)
	}
	return nil
}

type OutboxEntry struct {
	ID        int64           `db:"id"`
	Payload   json.RawMessage `db:"payload"`
	Status    OutboxStatus    `db:"status"`
	CreatedAt time.Time       `db:"created_at"`
}

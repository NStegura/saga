package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type EventStatus int

const (
	WAIT EventStatus = iota + 1
	DONE
)

func (os *EventStatus) String() string {
	return [...]string{"WAIT", "DONE"}[*os-1]
}

func (os *EventStatus) Value() (driver.Value, error) {
	return os.String(), nil
}

func (os *EventStatus) Scan(value any) error {
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("EventStatus should be a string")
	}

	switch strValue {
	case "WAIT":
		*os = WAIT
	case "DONE":
		*os = DONE
	default:
		return fmt.Errorf("invalid EventStatus: %s", strValue)
	}
	return nil
}

type EventEntry struct {
	ID        int64           `db:"id"`
	Payload   json.RawMessage `db:"payload"`
	Topic     string          `db:"topic"`
	Status    EventStatus     `db:"status"`
	CreatedAt time.Time       `db:"created_at"`
}

package models

import (
	"encoding/json"
)

type Event struct {
	ID      int64           `json:"id"`
	Payload json.RawMessage `json:"payload"`
	Topic   string          `json:"topic"`
}

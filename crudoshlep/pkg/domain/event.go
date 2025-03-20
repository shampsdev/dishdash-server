package domain

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID        string `gorm:"primary_key"`
	Name      string
	Data      json.RawMessage
	CreatedAt time.Time
}

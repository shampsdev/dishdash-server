package domain

import (
	"time"
)

type Lobby struct {
	ID        string
	CreatedAt time.Time
	Location  Coordinate
}

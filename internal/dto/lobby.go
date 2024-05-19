package dto

import "time"

type Lobby struct {
	ID        int64     `json:"id"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"createdAt"`
}

type LobbyToCreate struct {
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"createdAt"`
}

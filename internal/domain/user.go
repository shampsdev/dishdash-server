package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	Telegram  *int64    `json:"telegram"`
	CreatedAt time.Time `json:"createdAt"`
}

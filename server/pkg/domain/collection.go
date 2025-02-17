package domain

import "time"

type Collection struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Places      []*Place  `json:"places"`
	Visible     bool      `json:"visible"`
	Order       int64     `json:"order"`
}

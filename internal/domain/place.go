package domain

import (
	"time"
)

type Place struct {
	ID               int64      `json:"id"`
	Title            string     `json:"title"`
	ShortDescription string     `json:"shortDescription"`
	Description      string     `json:"description"`
	Images           []string   `json:"image"`
	Location         Coordinate `json:"location"`
	Address          string     `json:"address"`
	PriceMin         int        `json:"priceMin"`
	PriceMax         int        `json:"priceMax"`
	Tags             []*Tag     `json:"tags"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

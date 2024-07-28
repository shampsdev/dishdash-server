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
	PriceAvg         int        `json:"priceMin"`
	ReviewRating     float64    `json:"reviewRating"`
	ReviewCount      int        `json:"reviewCount"`
	Tags             []*Tag     `json:"tags"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

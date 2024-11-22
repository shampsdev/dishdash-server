package domain

import (
	"time"
)

type Place struct {
	ID               int64      `json:"id"`
	Title            string     `json:"title"`
	ShortDescription string     `json:"shortDescription"`
	Description      string     `json:"description"`
	Images           []string   `json:"images"`
	Location         Coordinate `json:"location"`
	Address          string     `json:"address"`
	PriceAvg         int        `json:"priceAvg"`
	ReviewRating     float64    `json:"reviewRating"`
	ReviewCount      int        `json:"reviewCount"`
	Tags             []*Tag     `json:"tags"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	Source           string     `json:"source"`
	Url              *string    `json:"url"`
}

func (place Place) Equals(other *Place) bool {
	return place.Title == other.Title && place.Address == other.Address
}

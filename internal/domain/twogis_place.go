package domain

import (
	"time"
)

type TwoGisPlace struct {
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	Lat          float64  `json:"lat"`
	Lon          float64  `json:"lon"`
	PhotoURL     string   `json:"photo_url"`
	ReviewRating float64  `json:"review_rating"`
	ReviewCount  int      `json:"review_count"`
	Rubrics      []string `json:"rubrics"`
	AveragePrice int      `json:"average_price"`
}

func (twoGisPlace *TwoGisPlace) ToPlace() *Place {
	return &Place{
		ID:               0,
		Title:            twoGisPlace.Name,
		ShortDescription: twoGisPlace.Address,
		Description:      twoGisPlace.Address,
		Images:           []string{twoGisPlace.PhotoURL},
		Location:         Coordinate{Lat: twoGisPlace.Lat, Lon: twoGisPlace.Lon},
		Address:          twoGisPlace.Address,
		PriceAvg:         twoGisPlace.AveragePrice,
		ReviewRating:     twoGisPlace.ReviewRating,
		ReviewCount:      twoGisPlace.ReviewCount,
		Tags:             nil, // apiTagRepo.SaveApiTag(ctx, &place)
		UpdatedAt:        time.Now(),
	}
}

package domain

import (
	"log"
	"strconv"
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

func (twoGisPlace *TwoGisPlace) parseTagToPlace() []*Tag {
	if twoGisPlace == nil {
		log.Println("twoGisPlace is nil")
		return nil
	}

	if len(twoGisPlace.Rubrics) == 0 {
		log.Println("twoGisPlace is empty")
		return nil
	}

	tags := make([]*Tag, len(twoGisPlace.Rubrics))
	for i, rubric := range twoGisPlace.Rubrics {
		if rubric == "" {
			log.Println("empty rubric found at index " + strconv.Itoa(i))
			return nil
		}
		tags[i] = &Tag{0, rubric, "no_icon"}
	}

	return tags
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
		Tags:             twoGisPlace.parseTagToPlace(),
		UpdatedAt:        time.Now(),
		Source:           "2gis",
	}
}

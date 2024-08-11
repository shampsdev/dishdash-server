package domain

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

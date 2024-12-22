package domain

type RecommendationType string

const (
	RecommendationTypeClassic     RecommendationType = "classic"
	RecommendationTypePriceBounds RecommendationType = "priceBound"
)

type RecommendationOpts struct {
	Type        RecommendationType            `json:"type"`
	Classic     *ClassicRecommendationOpts    `json:"classic"`
	PriceBounds *PriceBoundRecommendationOpts `json:"priceBound"`
}

type ClassicRecommendationOpts struct {
	PricePower float64 `json:"pricePower"`
	PriceCoeff float64 `json:"priceCoeff"`
	DistPower  float64 `json:"distPower"`
	DistCoeff  float64 `json:"distCoeff"`
}

type PriceBoundRecommendationOpts struct {
	PriceBound int `json:"priceBound"`
}

type RecommendData struct {
	Location Coordinate
	PriceAvg int
	Tags     []string
}

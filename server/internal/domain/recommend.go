package domain

type RecommendationType string

const (
	RecommendationTypeClassic RecommendationType = "classic"
)

type RecommendationOpts struct {
	Type    RecommendationType         `json:"type"`
	Classic *RecommendationOptsClassic `json:"classic"`
}

type RecommendationOptsClassic struct {
	PricePower float64 `json:"pricePower"`
	PriceCoeff float64 `json:"priceCoeff"`
	PriceBound int     `json:"priceBound"`
	DistPower  float64 `json:"distPower"`
	DistCoeff  float64 `json:"distCoeff"`
	DistBound  int     `json:"distBound"`
}

type RecommendData struct {
	Location Coordinate
	PriceAvg int
	Tags     []string
}

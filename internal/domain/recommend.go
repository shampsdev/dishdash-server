package domain

type RecommendOpts struct {
	PriceCoeff float64
	DistCoeff  float64
}

type RecommendData struct {
	Location Coordinate
	PriceAvg int
	Tags     []string
}

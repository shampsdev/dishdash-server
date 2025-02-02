package domain

type LobbySettings struct {
	Type             LobbyType                 `json:"type"`
	ClassicPlaces    *ClassicPlacesSettings    `json:"classicPlaces"`
	CollectionPlaces *CollectionPlacesSettings `json:"collectionPlaces"`
}

type ClassicPlacesSettings struct {
	Location       Coordinate          `json:"location"`
	Tags           []int64             `json:"tags"`
	PriceAvg       int                 `json:"priceAvg"`
	Recommendation *RecommendationOpts `json:"recommendation"`
}

type CollectionPlacesSettings struct {
	Location     *Coordinate `json:"location,omitempty"`
	CollectionID string      `json:"collectionId"`
}

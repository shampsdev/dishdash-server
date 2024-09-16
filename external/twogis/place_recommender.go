package twogis

import (
	"context"

	"dishdash.ru/internal/domain"
)

type PlaceRecommender struct{}

func (pr *PlaceRecommender) RecommendPlaces(
	_ context.Context,
	_ domain.RecommendOpts,
	data domain.RecommendData,
) ([]*domain.TwoGisPlace, error) {
	return FetchPlacesForLobbyFromAPI(data)
}

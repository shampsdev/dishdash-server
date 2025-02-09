package usecase

import (
	"context"
	"errors"
	"fmt"

	"dishdash.ru/cmd/server/config"
	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/repo"

	log "github.com/sirupsen/logrus"
)

type PlaceRecommender struct {
	dbPRRepo repo.PlaceRecommender

	pRepo repo.Place
	tRepo repo.Tag
}

func NewPlaceRecommender(
	dbPRRepo repo.PlaceRecommender,
	pRepo repo.Place,
	tRepo repo.Tag,
) *PlaceRecommender {
	return &PlaceRecommender{
		dbPRRepo: dbPRRepo,
		pRepo:    pRepo,
		tRepo:    tRepo,
	}
}

func (pr *PlaceRecommender) RecommendPlaces(
	ctx context.Context,
	settings domain.LobbySettings,
) ([]*domain.Place, error) {
	log.Debug("Starting recommendation process")
	var dbPlaces []*domain.Place
	var err error

	switch settings.Type {
	case domain.ClassicPlacesLobbyType:
		log.Debug("Using classic recommendation")
		if settings.ClassicPlaces == nil {
			return nil, errors.New("classic recommendation settings are chosen but not set")
		}

		if settings.ClassicPlaces.Recommendation == nil {
			settings.ClassicPlaces.Recommendation = defaultRecommendationOpts()
		}

		dbPlaces, err = pr.dbPRRepo.RecommendClassicPlaces(ctx, *settings.ClassicPlaces)
		if err != nil {
			return nil, fmt.Errorf("can't recommend from db: %w", err)
		}
		log.Debugf("Got %d places from db", len(dbPlaces))
	default:
		return nil, fmt.Errorf("unsupported recommendation type: %s", settings.Type)
	}

	return dbPlaces, nil
}

func defaultRecommendationOpts() *domain.RecommendationOpts {
	return &domain.RecommendationOpts{
		Type: domain.RecommendationTypeClassic,
		Classic: &domain.RecommendationOptsClassic{
			PricePower: config.C.Recommendation.PricePower,
			PriceCoeff: config.C.Recommendation.PriceCoeff,
			PriceBound: config.C.Recommendation.PriceBound,
			DistPower:  config.C.Recommendation.DistPower,
			DistCoeff:  config.C.Recommendation.DistCoeff,
			DistBound:  config.C.Recommendation.DistBound,
		},
	}
}

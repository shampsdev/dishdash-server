package usecase

import (
	"context"
	"errors"
	"fmt"

	"dishdash.ru/cmd/server/config"
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"

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
	opts *domain.RecommendationOpts,
	data domain.RecommendData,
) ([]*domain.Place, error) {
	log.Debug("Starting recommendation process")

	if opts == nil {
		opts = pr.defaultRecommendationOpts()
	}

	var dbPlaces []*domain.Place
	var err error

	switch opts.Type {
	case domain.RecommendationTypeClassic:
		log.Debug("Using classic recommendation")
		if opts.Classic == nil {
			return nil, errors.New("classic recommendation settings are chosen but not set")
		}

		dbPlaces, err = pr.dbPRRepo.RecommendClassic(ctx, *opts.Classic, data)
		if err != nil {
			return nil, fmt.Errorf("can't recommend from db: %w", err)
		}
		log.Debugf("Got %d places from db", len(dbPlaces))

	case domain.RecommendationTypePriceBounds:
		log.Debug("Using price bounds recommendation")
		if opts.PriceBounds == nil {
			return nil, errors.New("price bounds recommendation settings are chosen but not set")
		}

		dbPlaces, err = pr.dbPRRepo.RecommendPriceBound(ctx, *opts.PriceBounds, data)
		if err != nil {
			return nil, fmt.Errorf("can't recommend from db: %w", err)
		}
		log.Debugf("Got %d places from db", len(dbPlaces))
	default:
		return nil, fmt.Errorf("unknown recommendation type: %s", opts.Type)
	}

	if pr.goodEnough(dbPlaces, data) {
		return dbPlaces, nil
	}

	log.Debug("Can't recommend good enough places")

	return nil, errors.New("can't recommend places good enough")
}

func (pr *PlaceRecommender) defaultRecommendationOpts() *domain.RecommendationOpts {
	return &domain.RecommendationOpts{
		Type: domain.RecommendationTypeClassic,
		Classic: &domain.ClassicRecommendationOpts{
			PricePower: config.C.Recommendation.PricePower,
			PriceCoeff: config.C.Recommendation.PriceCoeff,
			DistPower:  config.C.Recommendation.DistPower,
			DistCoeff:  config.C.Recommendation.DistCoeff,
		},
	}
}

func (pr *PlaceRecommender) goodEnough(_ []*domain.Place, _ domain.RecommendData) bool {
	return true
}

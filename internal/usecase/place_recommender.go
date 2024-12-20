package usecase

import (
	"context"
	"errors"
	"fmt"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"

	log "github.com/sirupsen/logrus"
)

type PlaceRecommender struct {
	opts     domain.RecommendOpts
	dbPRRepo repo.PlaceRecommender

	pRepo repo.Place
	tRepo repo.Tag
}

func NewPlaceRecommender(
	opts domain.RecommendOpts,
	dbPRRepo repo.PlaceRecommender,
	pRepo repo.Place,
	tRepo repo.Tag,
) *PlaceRecommender {
	return &PlaceRecommender{
		opts:     opts,
		dbPRRepo: dbPRRepo,
		pRepo:    pRepo,
		tRepo:    tRepo,
	}
}

func (pr *PlaceRecommender) RecommendPlaces(ctx context.Context, data domain.RecommendData) ([]*domain.Place, error) {
	log.Debug("Starting recommendation process")

	dbPlaces, err := pr.dbPRRepo.RecommendPlaces(ctx, pr.opts, data)
	if err != nil {
		return nil, fmt.Errorf("can't recommend from db: %w", err)
	}
	log.Debugf("Got %d places from db", len(dbPlaces))

	if pr.goodEnough(dbPlaces, data) {
		return dbPlaces, nil
	}

	log.Debug("Can't recommend good enough places")

	return nil, errors.New("can't recommend places good enough")
}

func (pr *PlaceRecommender) goodEnough(_ []*domain.Place, _ domain.RecommendData) bool {
	return true
}

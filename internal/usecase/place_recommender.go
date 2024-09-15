package usecase

import (
	"context"
	"errors"
	"fmt"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"

	log "github.com/sirupsen/logrus"
)

type ApiPlaceRecommender interface {
	RecommendPlaces(
		ctx context.Context,
		opts domain.RecommendOpts,
		data domain.RecommendData,
	) ([]*domain.TwoGisPlace, error)
}

type PlaceSaver interface {
	SavePlace(ctx context.Context, placeInput SavePlaceInput) (*domain.Place, error)
}

type PlaceRecommender struct {
	opts         domain.RecommendOpts
	dbPRRepo     repo.PlaceRecommender
	twogisPRRepo ApiPlaceRecommender

	pRepo repo.Place
	tRepo repo.Tag
}

func NewPlaceRecommender(
	opts domain.RecommendOpts,
	dbPRRepo repo.PlaceRecommender,
	twogisPRRepo ApiPlaceRecommender,
	pRepo repo.Place,
	tRepo repo.Tag,
) *PlaceRecommender {
	return &PlaceRecommender{
		opts:         opts,
		dbPRRepo:     dbPRRepo,
		twogisPRRepo: twogisPRRepo,
		pRepo:        pRepo,
		tRepo:        tRepo,
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

	log.Debug("DB places not good enough, went to api")
	allTags, err := pr.tRepo.GetAllTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get all tags: %w", err)
	}

	twogisPlaces, err := pr.twogisPRRepo.RecommendPlaces(ctx, pr.opts, data)
	if err != nil {
		return nil, fmt.Errorf("can't get places from twogis: %w", err)
	}
	log.Debugf("Got %d places from twogis", len(twogisPlaces))

	for _, p := range twogisPlaces {
		err = pr.saveTwoGisPlace(ctx, p, allTags)
		if err != nil {
			return nil, fmt.Errorf("can't save place from twogis: %w", err)
		}
		log.Debugf("Save %s place from twogis", p.Name)
	}

	dbPlaces, err = pr.dbPRRepo.RecommendPlaces(ctx, pr.opts, data)
	if err != nil {
		return nil, fmt.Errorf("can't recommend from db: %w", err)
	}
	log.Debugf("Got %d places from db", len(dbPlaces))

	if pr.goodEnough(dbPlaces, data) {
		log.Debug("DB places good enough")
		return dbPlaces, nil
	}

	log.Debug("Can't recommend good enough places")

	return nil, errors.New("can't recommend places good enough")
}

func (pr *PlaceRecommender) goodEnough(_ []*domain.Place, _ domain.RecommendData) bool {
	return true
}

func (pr *PlaceRecommender) saveTwoGisPlace(ctx context.Context, p *domain.TwoGisPlace, allTags []*domain.Tag) error {
	id, err := pr.pRepo.SavePlace(ctx, p.ToPlace())
	if err != nil {
		return err
	}

	tagIDs := pr.twogisPlaceTagIDs(p, allTags)
	return pr.tRepo.AttachTagsToPlace(ctx, tagIDs, id)
}

func (pr *PlaceRecommender) twogisPlaceTagIDs(p *domain.TwoGisPlace, allTags []*domain.Tag) []int64 {
	tagIDs := make([]int64, 0)
	tagMap := make(map[string]int64, len(allTags))

	for _, tag := range allTags {
		tagMap[tag.Name] = tag.ID
	}
	for _, rubric := range p.Rubrics {
		tagID, found := tagMap[rubric]
		if found {
			tagIDs = append(tagIDs, tagID)
		}
	}

	return tagIDs
}

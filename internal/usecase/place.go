package usecase

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"dishdash.ru/external/twogis"

	"dishdash.ru/cmd/server/config"
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"
)

var (
	priceAvgLowerDelta = config.C.Defaults.PriceAvgLowerDelta
	priceAvgUpperDelta = config.C.Defaults.PriceAvgUpperDelta
)

type PlaceUseCase struct {
	tRepo repo.Tag
	pRepo repo.Place
}

func NewPlaceUseCase(tRepo repo.Tag, pRepo repo.Place) *PlaceUseCase {
	return &PlaceUseCase{tRepo: tRepo, pRepo: pRepo}
}

func (p PlaceUseCase) SavePlace(ctx context.Context, placeInput SavePlaceInput) (*domain.Place, error) {
	place := &domain.Place{
		Title:            placeInput.Title,
		ShortDescription: placeInput.ShortDescription,
		Description:      placeInput.Description,
		Images:           placeInput.Images,
		Location:         placeInput.Location,
		Address:          placeInput.Address,
		PriceAvg:         placeInput.PriceAvg,
		ReviewRating:     placeInput.ReviewRating,
		ReviewCount:      placeInput.ReviewCount,
	}
	id, err := p.pRepo.SavePlace(ctx, place)
	if err != nil {
		return nil, err
	}
	place.ID = id
	err = p.tRepo.AttachTagsToPlace(ctx, placeInput.Tags, id)
	if err != nil {
		return nil, err
	}

	place.Tags, err = p.tRepo.GetTagsByPlaceID(ctx, place.ID)
	if err != nil {
		return nil, err
	}

	return place, nil
}

func (p PlaceUseCase) SaveTwoGisPlace(ctx context.Context, twogisPlace *domain.TwoGisPlace) (int64, error) {
	placeId, err := p.pRepo.SaveTwoGisPlace(ctx, twogisPlace)
	if err != nil {
		return 0, err
	}
	return placeId, nil
}

func (p PlaceUseCase) GetPlaceByID(ctx context.Context, id int64) (*domain.Place, error) {
	place, err := p.pRepo.GetPlaceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	place.Tags, err = p.tRepo.GetTagsByPlaceID(ctx, id)
	if err != nil {
		return nil, err
	}
	return place, nil
}

func (p PlaceUseCase) GetAllPlaces(ctx context.Context) ([]*domain.Place, error) {
	places, err := p.pRepo.GetAllPlaces(ctx)
	if err != nil {
		return nil, err
	}
	for _, place := range places {
		place.Tags, err = p.tRepo.GetTagsByPlaceID(ctx, place.ID)
		if err != nil {
			return nil, err
		}
	}
	return places, nil
}

func getUniquePlaces(placesFromApi, placesFromBD []*domain.Place) []*domain.Place {
	// TODO: better asymptotic
	uniquePlaces := make([]*domain.Place, 0)
	uniquePlaces = append(uniquePlaces, placesFromApi...)

	for _, bdPlace := range placesFromBD {
		exists := false
		for _, apiPlace := range uniquePlaces {
			if bdPlace.Equals(apiPlace) {
				exists = true
				break
			}
		}
		if !exists {
			uniquePlaces = append(uniquePlaces, bdPlace)
		}
	}

	return uniquePlaces
}

func fetchTagIDsByNames(place *domain.TwoGisPlace, tags []*domain.Tag) ([]int64, error) {
	log.Debugf("Fetching tag IDs for api place: %s", place.Name)
	tagIDs := make([]int64, 0)
	tagMap := make(map[string]int64, len(tags))

	for _, tag := range tags {
		tagMap[tag.Name] = tag.ID
	}
	for _, rubric := range place.Rubrics {
		tagID, found := tagMap[rubric]
		if found {
			tagIDs = append(tagIDs, tagID)
		}
	}
	log.Debugf("Fetched %d tags for api place: %s", len(tagIDs), place.Name)
	return tagIDs, nil
}

func (p PlaceUseCase) GetPlacesForLobby(ctx context.Context, lobby *domain.Lobby) ([]*domain.Place, error) {
	log.Debugf("Starting GetPlacesForLobby for lobby ID: %s", lobby.ID)

	existingTags, err := p.tRepo.GetAllTags(ctx)
	if err != nil {
		log.WithError(err).Debugf("Failed to get tags from database for lobby ID: %s", lobby.ID)
		return nil, err
	}

	dbPlaces, err := p.pRepo.GetPlacesForLobby(ctx, lobby)
	if err != nil {
		log.WithError(err).Debugf("Failed to get places from database for lobby ID: %s", lobby.ID)
		return nil, err
	}

	for _, place := range dbPlaces {
		log.Infof("Fetching tags for place ID: %d", place.ID)
		place.Tags, err = p.tRepo.GetTagsByPlaceID(ctx, place.ID)
		if err != nil {
			log.WithError(err).Debugf("Failed to get tags for place ID: %d", place.ID)
			return nil, err
		}
	}

	if len(dbPlaces) < config.C.Defaults.MinDBPlaces {
		log.Debugf("Fewer than %d places found in DB for lobby ID: %s, fetching from 2GIS API.", config.C.Defaults.MinDBPlaces, lobby.ID)
		twoGisPlaces, err := twogis.FetchPlacesForLobbyFromAPI(lobby)
		if err != nil {
			log.WithError(err).Debugf("Failed to fetch places from 2GIS API for lobby ID: %s", lobby.ID)
			return nil, err
		}

		apiPlaces := make([]*domain.Place, 0)
		for _, twoGisPlace := range twoGisPlaces {
			log.Debugf("Processing 2GIS place: %s", twoGisPlace.Name)
			parsedPlace := twoGisPlace.ToPlace()

			// TODO: remove savings to DB, just attach
			placeTags, err := fetchTagIDsByNames(twoGisPlace, existingTags)
			if err != nil {
				log.WithError(err).Errorf("Failed to fetch tag IDs for 2GIS place: %s", twoGisPlace.Name)
				return nil, err
			}

			placeId, err := p.SaveTwoGisPlace(ctx, twoGisPlace)
			if errors.Is(err, repo.ErrPlaceExists) {
				log.Debugf("Place already exists in DB, skipping 2GIS place: %s", twoGisPlace.Name)
				continue
			}
			if err != nil {
				log.WithError(err).Errorf("Failed to save 2GIS place: %s", twoGisPlace.Name)
				return nil, err
			}

			parsedPlace.ID = placeId
			apiPlaces = append(apiPlaces, parsedPlace)

			err = p.tRepo.AttachTagsToPlace(ctx, placeTags, placeId)
			if err != nil {
				log.WithError(err).Errorf("Failed to attach tags to place ID: %d", placeId)
				return nil, err
			}
		}

		filteredPlaces := make([]*domain.Place, 0)
		for _, place := range apiPlaces {
			if place.PriceAvg > lobby.PriceAvg-priceAvgLowerDelta && place.PriceAvg < lobby.PriceAvg+priceAvgUpperDelta {
				filteredPlaces = append(filteredPlaces, place)
			}
		}
		log.Debugf("Found filtered %d places for lobby ID: %s", len(filteredPlaces), lobby.ID)

		log.Debugf("Returning filtered unique places for lobby ID: %s", lobby.ID)
		return getUniquePlaces(filteredPlaces, dbPlaces), nil
	}

	log.Debugf("Returning DB places for lobby ID: %s", lobby.ID)
	return dbPlaces, nil
}

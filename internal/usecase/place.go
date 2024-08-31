package usecase

import (
	"context"
	"errors"
	"log"

	"dishdash.ru/external/twogis"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"
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

func (p PlaceUseCase) GetPlacesForLobby(ctx context.Context, lobby *domain.Lobby) ([]*domain.Place, error) {
	log.Printf("[INFO] Starting GetPlacesForLobby for lobby ID: %d", lobby.ID)

	dbPlaces, err := p.pRepo.GetPlacesForLobby(ctx, lobby)
	if err != nil {
		log.Printf("[ERROR] Failed to get places from database for lobby ID: %d, error: %v", lobby.ID, err)
		return nil, err
	}

	for _, place := range dbPlaces {
		log.Printf("[INFO] Fetching tags for place ID: %d", place.ID)
		place.Tags, err = p.tRepo.GetTagsByPlaceID(ctx, place.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to get tags for place ID: %d, error: %v", place.ID, err)
			return nil, err
		}
	}

	if len(dbPlaces) <= 5 {
		log.Printf("[INFO] Fewer than 5 places found in DB for lobby ID: %d, fetching from 2GIS API.", lobby.ID)
		twoGisPlaces, err := twogis.FetchPlacesForLobbyFromAPI(lobby)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch places from 2GIS API for lobby ID: %d, error: %v", lobby.ID, err)
			return nil, err
		}

		apiPlaces := make([]*domain.Place, 0)
		for _, twoGisPlace := range twoGisPlaces {
			log.Printf("[INFO] Processing 2GIS place: %s", twoGisPlace.Name)
			parsedPlace := twoGisPlace.ToPlace()
			
			tags, err := p.tRepo.SaveApiTag(ctx, twoGisPlace)
			if err != nil {
				log.Printf("[ERROR] Failed to save tags for 2GIS place: %s, error: %v", twoGisPlace.Name, err)
				return nil, err
			}
			
			placeId, err := p.SaveTwoGisPlace(ctx, twoGisPlace)
			if errors.Is(err, repo.ErrPlaceExists) {
				log.Printf("[INFO] Place already exists in DB, skipping 2GIS place: %s", twoGisPlace.Name)
				continue
			}
			if err != nil {
				log.Printf("[ERROR] Failed to save 2GIS place: %s, error: %v", twoGisPlace.Name, err)
				return nil, err
			}
			
			parsedPlace.ID = placeId
			apiPlaces = append(apiPlaces, parsedPlace)

			err = p.tRepo.AttachTagsToPlace(ctx, tags, placeId)
			if err != nil {
				log.Printf("[ERROR] Failed to attach tags to place ID: %d, error: %v", placeId, err)
				return nil, err
			}
		}

		log.Printf("[INFO] Returning unique places for lobby ID: %d", lobby.ID)
		return getUniquePlaces(apiPlaces, dbPlaces), nil
	}

	log.Printf("[INFO] Returning DB places for lobby ID: %d", lobby.ID)
	return dbPlaces, nil
}


package usecase

import (
	"context"
	"errors"

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
	dbPlaces, err := p.pRepo.GetPlacesForLobby(ctx, lobby)
	for _, place := range dbPlaces {
		place.Tags, err = p.tRepo.GetTagsByPlaceID(ctx, place.ID)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	if len(dbPlaces) <= 5 {
		twoGisPlaces, err := twogis.FetchPlacesForLobbyFromAPI(lobby)
		if err != nil {
			return nil, err
		}
		apiPlaces := make([]*domain.Place, 0)
		for _, twoGisPlace := range twoGisPlaces {
			parsedPlace := twoGisPlace.ToPlace()
			tags, err := p.tRepo.SaveApiTag(ctx, twoGisPlace)
			if err != nil {
				return nil, err
			}
			placeId, err := p.SaveTwoGisPlace(ctx, twoGisPlace)
			if errors.Is(err, repo.ErrPlaceExists) {
				continue
			}
			parsedPlace.ID = placeId
			apiPlaces = append(apiPlaces, parsedPlace)
			if err != nil {
				return nil, err
			}
			err = p.tRepo.AttachTagsToPlace(ctx, tags, placeId)
			if err != nil {
				return nil, err
			}
		}

		return getUniquePlaces(apiPlaces, dbPlaces), nil
	}

	return dbPlaces, nil
}

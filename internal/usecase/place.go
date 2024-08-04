package usecase

import (
	"context"

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

	return place, nil
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

func (p PlaceUseCase) GetPlacesForLobby(ctx context.Context, lobby *domain.Lobby) ([]*domain.Place, error) {
	places, err := p.pRepo.GetPlacesForLobby(ctx, lobby)
	if err != nil {
		return nil, err
	}

	return places, nil
}

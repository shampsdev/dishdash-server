package usecase

import (
	"context"

	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/repo"
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
		Source:           placeInput.Source,
		Url:              placeInput.Url,
		Boost:            placeInput.Boost,
		BoostRadius:      placeInput.BoostRadius,
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

func (p PlaceUseCase) UpdatePlace(ctx context.Context, placeInput UpdatePlaceInput) (*domain.Place, error) {
	place := &domain.Place{
		ID:               placeInput.ID,
		Title:            placeInput.Title,
		ShortDescription: placeInput.ShortDescription,
		Description:      placeInput.Description,
		Images:           placeInput.Images,
		Location:         placeInput.Location,
		Address:          placeInput.Address,
		PriceAvg:         placeInput.PriceAvg,
		ReviewRating:     placeInput.ReviewRating,
		ReviewCount:      placeInput.ReviewCount,
		Source:           placeInput.Source,
		Url:              placeInput.Url,
		Boost:            placeInput.Boost,
		BoostRadius:      placeInput.BoostRadius,
	}
	err := p.pRepo.UpdatePlace(ctx, place)
	if err != nil {
		return nil, err
	}
	err = p.tRepo.DetachTagsFromPlace(ctx, place.ID)
	if err != nil {
		return nil, err
	}
	err = p.tRepo.AttachTagsToPlace(ctx, placeInput.Tags, place.ID)
	if err != nil {
		return nil, err
	}

	place.Tags, err = p.tRepo.GetTagsByPlaceID(ctx, placeInput.ID)
	if err != nil {
		return nil, err
	}

	return place, nil
}

func (p PlaceUseCase) DeletePlace(ctx context.Context, id int64) error {
	err := p.pRepo.DeletePlace(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (p PlaceUseCase) GetPlaceByID(ctx context.Context, id int64) (*domain.Place, error) {
	place, err := p.pRepo.GetPlaceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return place, nil
}

func (p PlaceUseCase) GetPlaceByUrl(ctx context.Context, url string) (*domain.Place, error) {
	place, err := p.pRepo.GetPlaceByUrl(ctx, url)
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
	return places, nil
}

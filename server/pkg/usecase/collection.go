package usecase

import (
	"context"

	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/repo"
)

type CollectionUseCase struct {
	cRepo repo.Collection
}


func NewCollectionUseCase(cRepo repo.Collection) *CollectionUseCase {
    return &CollectionUseCase{cRepo: cRepo}
}


func (cu *CollectionUseCase) SaveCollection(ctx context.Context, saveCollectionInput SaveCollectionInput) (*domain.Collection, error) {
    collection := &domain.Collection{
        Name:        saveCollectionInput.Name,
        Description: saveCollectionInput.Description,
        Avatar:      saveCollectionInput.Avatar,
        Visible:     saveCollectionInput.Visible,
        Order:       saveCollectionInput.Order,
    }

    id, err := cu.cRepo.SaveCollection(ctx, collection)
    if err != nil {
        return nil, err
    }

    collection.ID = id
    err = cu.cRepo.AttachPlacesToCollection(ctx, saveCollectionInput.Places, id)
    if err != nil {
        return nil, err
    }
    return collection, nil
}

func (cu *CollectionUseCase) GetAllCollections(ctx context.Context) ([]*domain.Collection, error) {
    return cu.cRepo.GetAllCollections(ctx)
}

func (cu *CollectionUseCase) DeleteCollection(ctx context.Context, collectionID int64) error {
    return cu.cRepo.DeleteCollectionByID(ctx, collectionID)
}




func (cu *CollectionUseCase) UpdateCollection(ctx context.Context, updateCollectionInput UpdateCollectionInput) (*domain.Collection, error) {
    collection := &domain.Collection{
        ID:          updateCollectionInput.ID,
        Name:        updateCollectionInput.Name,
        Description: updateCollectionInput.Description,
        Avatar:      updateCollectionInput.Avatar,
        Visible:     updateCollectionInput.Visible,
        Order:       updateCollectionInput.Order,
    }

    err := cu.cRepo.UpdateCollection(ctx, collection)
    if err != nil {
        return nil, err
    }

    err = cu.cRepo.DetachPlacesFromCollection(ctx, updateCollectionInput.ID)
    if err != nil {
        return nil, err
    }
    err = cu.cRepo.AttachPlacesToCollection(ctx, updateCollectionInput.Places, updateCollectionInput.ID)
    if err != nil {
        return nil, err
    }


    return collection, nil
}
package usecase

import (
	"context"

	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/repo"
)

type TagUseCase struct {
	tRepo repo.Tag
}

func NewTagUseCase(tRepo repo.Tag) *TagUseCase {
	return &TagUseCase{tRepo: tRepo}
}

func (t TagUseCase) SaveTag(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	id, err := t.tRepo.SaveTag(ctx, tag)
	if err != nil {
		return nil, err
	}
	tag.ID = id
	return tag, nil
}

func (t TagUseCase) GetAllTags(ctx context.Context) ([]*domain.Tag, error) {
	return t.tRepo.GetAllTags(ctx)
}

func (t TagUseCase) DeleteTag(ctx context.Context, tagId int64) error {
	return t.tRepo.DeleteTag(ctx, tagId)
}

func (t TagUseCase) UpdateTag(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	return t.tRepo.UpdateTag(ctx, tag)
}

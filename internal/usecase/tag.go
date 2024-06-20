package usecase

import (
	"context"

	"dishdash.ru/internal/repo"

	"dishdash.ru/internal/domain"
)

type TagUseCase struct {
	tagRepo repo.Tag
}

func NewTagUseCase(tagRepo repo.Tag) *TagUseCase {
	return &TagUseCase{tagRepo: tagRepo}
}

func (t *TagUseCase) CreateTag(ctx context.Context, tagInput TagInput) (*domain.Tag, error) {
	tag := &domain.Tag{
		Name: tagInput.Name,
		Icon: tagInput.Icon,
	}
	id, err := t.tagRepo.CreateTag(ctx, tag)
	if err != nil {
		return nil, err
	}
	tag.ID = id
	return tag, nil
}

func (t *TagUseCase) GetAllTags(ctx context.Context) ([]*domain.Tag, error) {
	return t.tagRepo.GetAllTags(ctx)
}

package usecase

import (
	"context"

	"dishdash.ru/internal/repo"

	"dishdash.ru/internal/domain"
)

type Tag struct {
	tagRepo repo.Tag
}

func NewTag(tagRepo repo.Tag) *Tag {
	return &Tag{tagRepo: tagRepo}
}

func (t *Tag) CreateTag(ctx context.Context, tag *domain.Tag) error {
	id, err := t.tagRepo.CreateTag(ctx, tag)
	if err != nil {
		return err
	}
	tag.ID = id
	return nil
}

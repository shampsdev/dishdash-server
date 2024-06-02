package usecase

import (
	"context"
	"dishdash.ru/internal/domain"
)

type Tag struct {
	tagRepo TagRepository
}

func NewTag(tagRepo TagRepository) *Tag {
	return &Tag{tagRepo: tagRepo}
}

func (t *Tag) SaveTag(ctx context.Context, tag *domain.Tag) error {
	return t.tagRepo.SaveTag(ctx, tag)
}

func (t *Tag) AttachTagToCard(ctx context.Context, tagID int64, cardID int64) error {
	return t.tagRepo.AttachTagToCard(ctx, tagID, cardID)
}

func (t *Tag) GetTagsByCardID(ctx context.Context, cardID int64) ([]*domain.Tag, error) {
	return t.tagRepo.GetTagsByCardID(ctx, cardID)
}

package repo

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Card interface {
	CreateCard(ctx context.Context, card *domain.Card) (int64, error)
	GetCardByID(ctx context.Context, id int64) (*domain.Card, error)
	GetAllCards(ctx context.Context) ([]*domain.Card, error)
}

type Tag interface {
	CreateTag(ctx context.Context, tag *domain.Tag) (int64, error)
	AttachTagsToCard(ctx context.Context, tagIDs []int64, cardID int64) error
	GetTagsByCardID(ctx context.Context, cardID int64) ([]*domain.Tag, error)
}

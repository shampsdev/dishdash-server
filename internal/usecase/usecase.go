package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Cases struct {
	Card Card
	Tag  Tag
}

type CardInput struct {
	Title            string            `json:"title"`
	ShortDescription string            `json:"shortDescription"`
	Description      string            `json:"description"`
	Image            string            `json:"image"`
	Location         domain.Coordinate `json:"location"`
	Address          string            `json:"address"`
	Price            int               `json:"price"`
	Tags             []int64           `json:"tags"`
}

type Card interface {
	CreateCard(ctx context.Context, card CardInput) (*domain.Card, error)
	GetCardByID(ctx context.Context, id int64) (*domain.Card, error)
	GetAllCards(ctx context.Context) ([]*domain.Card, error)
}

type TagInput struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Tag interface {
	CreateTag(ctx context.Context, tag TagInput) (*domain.Tag, error)
}

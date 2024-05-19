package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Card struct {
	cardRepo CardRepository
}

func NewCard(cardRepo CardRepository) *Card {
	return &Card{cardRepo: cardRepo}
}

func (c *Card) SaveCard(ctx context.Context, card *domain.Card) error {
	return c.cardRepo.SaveCard(ctx, card)
}

func (c *Card) GetCards(ctx context.Context) ([]*domain.Card, error) {
	return c.cardRepo.GetCards(ctx)
}

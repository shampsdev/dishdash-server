package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Card struct {
	cardRepo CardRepository
	tagRepo  TagRepository
}

func NewCard(cardRepo CardRepository, tagRepo TagRepository) *Card {
	return &Card{cardRepo: cardRepo, tagRepo: tagRepo}
}

func (c *Card) SaveCard(ctx context.Context, card *domain.Card) error {
	return c.cardRepo.SaveCard(ctx, card)
}

func (c *Card) GetCards(ctx context.Context) ([]*domain.Card, error) {
	cards, err := c.cardRepo.GetCards(ctx)
	if err != nil {
		return nil, err
	}

	for _, card := range cards {
		tags, err := c.tagRepo.GetTagsByCardID(ctx, card.ID)
		if err != nil {
			return nil, err
		}
		card.Tags = tags
	}

	return cards, nil
}

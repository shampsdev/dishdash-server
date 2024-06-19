package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"
)

type Card struct {
	cardRepo repo.Card
	tagRepo  repo.Tag
}

func NewCard(cardRepo repo.Card, tagRepo repo.Tag) *Card {
	return &Card{cardRepo: cardRepo, tagRepo: tagRepo}
}

func (c *Card) CreateCard(ctx context.Context, card *domain.Card) error {
	id, err := c.cardRepo.CreateCard(ctx, card)
	if err != nil {
		return err
	}
	card.ID = id

	return nil
}

func (c *Card) GetCardByID(ctx context.Context, id int64) (*domain.Card, error) {
	card, err := c.cardRepo.GetCardByID(ctx, id)
	if err != nil {
		return nil, err
	}
	card.Tags, err = c.tagRepo.GetTagsByCardID(ctx, card.ID)
	return card, err
}

func (c *Card) AttachTagsToCard(ctx context.Context, tagIDs []int64, cardID int64) error {
	return c.tagRepo.AttachTagsToCard(ctx, tagIDs, cardID)
}

func (c *Card) GetAllCards(ctx context.Context) ([]*domain.Card, error) {
	cards, err := c.cardRepo.GetAllCards(ctx)
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

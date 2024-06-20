package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"
)

type CardUseCase struct {
	cardRepo repo.Card
	tagRepo  repo.Tag
}

func NewCardUseCase(cardRepo repo.Card, tagRepo repo.Tag) *CardUseCase {
	return &CardUseCase{cardRepo: cardRepo, tagRepo: tagRepo}
}

func (c *CardUseCase) CreateCard(ctx context.Context, cardInput CardInput) (*domain.Card, error) {
	card := &domain.Card{
		Title:            cardInput.Title,
		ShortDescription: cardInput.ShortDescription,
		Description:      cardInput.Description,
		Image:            cardInput.Image,
		Location:         cardInput.Location,
		Address:          cardInput.Address,
		PriceMin:         cardInput.PriceMin,
		PriceMax:         cardInput.PriceMax,
		Tags:             nil,
	}
	id, err := c.cardRepo.CreateCard(ctx, card)
	if err != nil {
		return nil, err
	}
	card.ID = id
	err = c.tagRepo.AttachTagsToCard(ctx, cardInput.Tags, id)
	if err != nil {
		return nil, err
	}

	return c.cardRepo.GetCardByID(ctx, id)
}

func (c *CardUseCase) GetCardByID(ctx context.Context, id int64) (*domain.Card, error) {
	card, err := c.cardRepo.GetCardByID(ctx, id)
	if err != nil {
		return nil, err
	}
	card.Tags, err = c.tagRepo.GetTagsByCardID(ctx, card.ID)
	return card, err
}

func (c *CardUseCase) AttachTagsToCard(ctx context.Context, tagIDs []int64, cardID int64) error {
	return c.tagRepo.AttachTagsToCard(ctx, tagIDs, cardID)
}

func (c *CardUseCase) GetAllCards(ctx context.Context) ([]*domain.Card, error) {
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

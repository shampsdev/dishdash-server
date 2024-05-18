package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
)

type CardRepository interface {
	SaveCard(ctx context.Context, card *domain.Card) error
	GetCards(ctx context.Context) ([]domain.Card, error)
}

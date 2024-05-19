package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Cases struct {
	Card  *Card
	Lobby *Lobby
	Swipe *Swipe
}

type CardRepository interface {
	SaveCard(ctx context.Context, card *domain.Card) error
	GetCards(ctx context.Context) ([]*domain.Card, error)
}

type LobbyRepository interface {
	GetLobbyByID(ctx context.Context, id int64) (*domain.Lobby, error)
	SaveLobby(ctx context.Context, lobby *domain.Lobby) error
}

type SwipeRepository interface {
	SaveSwipe(ctx context.Context, swipe *domain.Swipe) error
}

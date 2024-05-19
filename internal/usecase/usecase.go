package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Cases struct {
	Card  *Card
	Lobby *Lobby
}

type CardRepository interface {
	SaveCard(ctx context.Context, card *domain.Card) error
	GetCards(ctx context.Context) ([]*domain.Card, error)
}

type LobbyRepository interface {
	GetLobbyByID(ctx context.Context, id int64) (*domain.Lobby, error)
	SaveLobby(ctx context.Context, lobby *domain.Lobby) error
}

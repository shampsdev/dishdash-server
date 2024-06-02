package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Cases struct {
	Card  *Card
	Lobby *Lobby
	Swipe *Swipe
	Tag   *Tag
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

type TagRepository interface {
	SaveTag(ctx context.Context, tag *domain.Tag) error
	AttachTagToCard(ctx context.Context, tagID, cardID int64) error
	GetTagsByCardID(ctx context.Context, cardID int64) ([]*domain.Tag, error)
}

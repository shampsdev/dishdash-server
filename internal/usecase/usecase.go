package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Cases struct {
	Card  Card
	Tag   Tag
	Lobby Lobby
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
	CreateCard(ctx context.Context, cardInput CardInput) (*domain.Card, error)
	GetCardByID(ctx context.Context, id int64) (*domain.Card, error)
	GetAllCards(ctx context.Context) ([]*domain.Card, error)
}

type TagInput struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Tag interface {
	CreateTag(ctx context.Context, tagInput TagInput) (*domain.Tag, error)
}

type LobbyInput struct {
	Location domain.Coordinate `json:"location"`
}

type Lobby interface {
	CreateLobby(ctx context.Context, lobbyInput LobbyInput) (*domain.Lobby, error)
	DeleteLobbyByID(ctx context.Context, id string) error
	NearestLobby(ctx context.Context, loc domain.Coordinate) (*domain.Lobby, float64, error)
}

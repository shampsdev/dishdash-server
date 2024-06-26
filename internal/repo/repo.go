package repo

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Card interface {
	CreateCard(ctx context.Context, card *domain.Card) (int64, error)
	GetCardByID(ctx context.Context, id int64) (*domain.Card, error)
	GetAllCards(ctx context.Context) ([]*domain.Card, error)
}

type Tag interface {
	CreateTag(ctx context.Context, tag *domain.Tag) (int64, error)
	AttachTagsToCard(ctx context.Context, tagIDs []int64, cardID int64) error
	GetTagsByCardID(ctx context.Context, cardID int64) ([]*domain.Tag, error)
	GetAllTags(ctx context.Context) ([]*domain.Tag, error)
}

type Lobby interface {
	CreateLobby(ctx context.Context, lobby *domain.Lobby) (*domain.Lobby, error)
	DeleteLobbyByID(ctx context.Context, lobbyID string) error
	NearestLobby(ctx context.Context, loc domain.Coordinate) (lobby *domain.Lobby, dist float64, err error)
	GetLobbyByID(ctx context.Context, id string) (*domain.Lobby, error)
}

type User interface {
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	UpdateUser(ctx context.Context, user *domain.User) (string, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
}

type Swipe interface {
	CreateSwipe(ctx context.Context, swipe *domain.Swipe) error
}

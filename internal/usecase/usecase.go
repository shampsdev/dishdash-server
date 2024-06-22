package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Cases struct {
	Card  Card
	Tag   Tag
	Lobby Lobby
	User  User
}

type CardInput struct {
	Title            string            `json:"title"`
	ShortDescription string            `json:"shortDescription"`
	Description      string            `json:"description"`
	Image            string            `json:"image"`
	Location         domain.Coordinate `json:"location"`
	Address          string            `json:"address"`
	PriceMin         int               `json:"priceMin"`
	PriceMax         int               `json:"priceMax"`
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
	GetAllTags(ctx context.Context) ([]*domain.Tag, error)
}

type LobbyInput struct {
	Location domain.Coordinate `json:"location"`
}

type LobbySettingsInput struct {
	LobbyID     string  `json:"lobbyID"`
	PriceMin    int     `json:"priceMin"`
	PriceMax    int     `json:"priceMax"`
	MaxDistance float64 `json:"maxDistance"`
	Tags        []*Tag  `json:"tags"`
}

type Lobby interface {
	CreateLobby(ctx context.Context, lobbyInput LobbyInput) (*domain.Lobby, error)
	DeleteLobbyByID(ctx context.Context, id string) error
	NearestLobby(ctx context.Context, loc domain.Coordinate) (*domain.Lobby, float64, error)
	//GetLobbyByID(ctx context.Context, id string) (*domain.Lobby, error)
	//
	//ChangeLobbySettings(ctx context.Context, lobbyID string, settings LobbySettingsInput) (*domain.LobbySettings, error)
	//FillLobbyWithCards(ctx context.Context, lobbyID string) ([]*domain.Card, error)
	//SwipeCard(ctx context.Context, lobbyID string, userID string, cardID int) (*domain.Match, error)
	//VoteForCard(ctx context.Context, lobbyID string, userID string, cardID int) error
}

type UserInput struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type User interface {
	CreateUser(ctx context.Context, userInput UserInput) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
}

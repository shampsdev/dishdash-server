package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Cases struct {
	Tag      Tag
	User     User
	Place    Place
	Swipe    Swipe
	Lobby    Lobby
	RoomRepo RoomRepo
}

type Tag interface {
	SaveTag(ctx context.Context, tag *domain.Tag) (*domain.Tag, error)
	GetAllTags(ctx context.Context) ([]*domain.Tag, error)
	SaveApiTag(ctx context.Context, place *domain.TwoGisPlace) ([]int64, error)
}

type User interface {
	SaveUser(ctx context.Context, user *domain.User) (*domain.User, error)
	SaveUserWithID(ctx context.Context, user *domain.User, id string) error
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)

	AttachUserToLobby(ctx context.Context, userID, lobbyID string) error
}

type SavePlaceInput struct {
	Title            string            `json:"title"`
	ShortDescription string            `json:"shortDescription"`
	Description      string            `json:"description"`
	Location         domain.Coordinate `json:"location"`
	Address          string            `json:"address"`
	PriceAvg         int               `json:"priceMin"`
	ReviewRating     float64           `json:"reviewRating"`
	ReviewCount      int               `json:"reviewCount"`
	Images           []string          `json:"images"`
	Tags             []int64           `json:"tags"`
}

type Place interface {
	SavePlace(ctx context.Context, placeInput SavePlaceInput) (*domain.Place, error)
	SaveTwoGisPlace(ctx context.Context, twogisPlace *domain.TwoGisPlace) (int64, error)
	GetPlaceByID(ctx context.Context, id int64) (*domain.Place, error)
	// GetAllPlaces is very long operation now
	GetAllPlaces(ctx context.Context) ([]*domain.Place, error)
	GetPlacesForLobby(ctx context.Context, lobby *domain.Lobby) ([]*domain.Place, error)
}

type SaveLobbyInput struct {
	PriceAvg int               `json:"priceAvg"`
	Location domain.Coordinate `json:"location"`
	Tags     []int64           `json:"tags"`
	Places   []int64           `json:"places"`
}

type UpdateLobbyInput struct {
	ID string
	SaveLobbyInput
}

type FindLobbyInput struct {
	Dist     float64           `json:"dist"`
	Location domain.Coordinate `json:"location"`
}

type Lobby interface {
	SaveLobby(ctx context.Context, lobbyInput SaveLobbyInput) (*domain.Lobby, error)
	UpdateLobby(ctx context.Context, lobbyInput UpdateLobbyInput) (*domain.Lobby, error)
	DeleteLobbyByID(ctx context.Context, id string) error
	GetLobbyByID(ctx context.Context, id string) (*domain.Lobby, error)

	FindLobby(ctx context.Context, input FindLobbyInput) (*domain.Lobby, error)
	NearestActiveLobby(ctx context.Context, loc domain.Coordinate) (*domain.Lobby, float64, error)
}

type Swipe interface {
	SaveSwipe(ctx context.Context, swipe *domain.Swipe) error
}

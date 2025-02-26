package usecase

import (
	"context"

	"dishdash.ru/pkg/domain"
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
	DeleteTag(ctx context.Context, tagId int64) error
	UpdateTag(ctx context.Context, tag *domain.Tag) (*domain.Tag, error)
}

type User interface {
	SaveUser(ctx context.Context, user *domain.User) (*domain.User, error)
	SaveUserWithID(ctx context.Context, user *domain.User, id string) error
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByTelegram(ctx context.Context, telegram *int64) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	GetUsersByLobbyID(ctx context.Context, lobbyID string) ([]*domain.User, error)
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
	Source           string            `json:"source"`
	Url              *string           `json:"url"`
	Boost            *float64          `json:"boost"`
	BoostRadius      *float64          `json:"boostRadius"`
	Images           []string          `json:"images"`
	Tags             []int64           `json:"tags"`
}

type UpdatePlaceInput struct {
	ID int64
	SavePlaceInput
}

type Place interface {
	SavePlace(ctx context.Context, placeInput SavePlaceInput) (*domain.Place, error)
	UpdatePlace(ctx context.Context, place UpdatePlaceInput) (*domain.Place, error)
	DeletePlace(ctx context.Context, id int64) error
	GetPlaceByID(ctx context.Context, id int64) (*domain.Place, error)
	GetPlaceByUrl(ctx context.Context, url string) (*domain.Place, error)
	// GetAllPlaces is very long operation now
	GetAllPlaces(ctx context.Context) ([]*domain.Place, error)
}

type Lobby interface {
	CreateLobby(ctx context.Context, settings domain.LobbySettings) (*domain.Lobby, error)
	DeleteLobbyByID(ctx context.Context, id string) error
	GetLobbyByID(ctx context.Context, id string) (*domain.Lobby, error)

	SetLobbySettings(ctx context.Context, lobbyID string, settings domain.LobbySettings) error
	SetLobbyState(ctx context.Context, lobbyID string, state domain.LobbyState) error
	SetLobbyUsers(ctx context.Context, lobbyID string, userIDs []string) ([]*domain.User, error)
	AttachOrderedPlacesToLobby(ctx context.Context, placeIDs []int64, lobbyID string) error
}

type Swipe interface {
	SaveSwipe(ctx context.Context, swipe *domain.Swipe) error
	GetCount(ctx context.Context) (int, error)
	GetSwipesByLobbyID(ctx context.Context, lobbyID string) ([]*domain.Swipe, error)
}

type Collection interface {
	SaveCollection(ctx context.Context, collection *domain.Collection) (*domain.Collection, error)
	GetAllCollections(ctx context.Context) ([]*domain.Collection, error)
	DeleteCollection(ctx context.Context, collectionID int64) error
	UpdateCollection(ctx context.Context, tag *domain.Tag) (*domain.Collection, error)
}

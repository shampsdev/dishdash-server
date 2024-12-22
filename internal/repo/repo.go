package repo

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Tag interface {
	SaveTag(ctx context.Context, tag *domain.Tag) (int64, error)
	GetAllTags(ctx context.Context) ([]*domain.Tag, error)

	DeleteTag(ctx context.Context, tagId int64) error
	UpdateTag(ctx context.Context, tag *domain.Tag) (*domain.Tag, error)

	AttachTagsToPlace(ctx context.Context, tagIDs []int64, placeID int64) error
	DetachTagsFromPlace(ctx context.Context, placeID int64) error
	GetTagsByPlaceID(ctx context.Context, placeID int64) ([]*domain.Tag, error)

	AttachTagsToLobby(ctx context.Context, tagIDs []int64, lobbyID string) error
	DetachTagsFromLobby(ctx context.Context, lobbyID string) error
	GetTagsByLobbyID(ctx context.Context, lobbyID string) ([]*domain.Tag, error)
}

type Place interface {
	SavePlace(ctx context.Context, place *domain.Place) (int64, error)
	UpdatePlace(ctx context.Context, place *domain.Place) error
	DeletePlace(ctx context.Context, id int64) error
	GetPlaceByID(ctx context.Context, id int64) (*domain.Place, error)
	GetPlaceByUrl(ctx context.Context, url string) (*domain.Place, error)
	GetAllPlaces(ctx context.Context) ([]*domain.Place, error)

	DetachPlacesFromLobby(ctx context.Context, lobbyID string) error
	AttachPlacesToLobby(ctx context.Context, placeIDs []int64, lobbyID string) error
	GetPlacesByLobbyID(ctx context.Context, lobbyID string) ([]*domain.Place, error)
	GetPlacesForLobby(ctx context.Context, lobby *domain.Lobby) ([]*domain.Place, error)
}

type PlaceRecommender interface {
	RecommendClassic(ctx context.Context, opts domain.ClassicRecommendationOpts, data domain.RecommendData) ([]*domain.Place, error)
	RecommendPriceBound(ctx context.Context, opts domain.PriceBoundRecommendationOpts, data domain.RecommendData) ([]*domain.Place, error)
}

type User interface {
	SaveUser(ctx context.Context, user *domain.User) (string, error)
	SaveUserWithID(ctx context.Context, user *domain.User, id string) error
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByTelegram(ctx context.Context, telegram *int64) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)

	AttachUsersToLobby(ctx context.Context, userID []string, lobbyID string) error
	DetachUsersFromLobby(ctx context.Context, lobbyID string) error
	GetUsersByLobbyID(ctx context.Context, lobbyID string) ([]*domain.User, error)
}

type Lobby interface {
	SaveLobby(ctx context.Context, lobby *domain.Lobby) (string, error)
	DeleteLobbyByID(ctx context.Context, id string) error
	GetLobbyByID(ctx context.Context, id string) (*domain.Lobby, error)

	UpdateLobby(ctx context.Context, lobby *domain.Lobby) error
	SetLobbyState(ctx context.Context, lobbyID string, state domain.LobbyState) error

	NearestActiveLobbyID(ctx context.Context, loc domain.Coordinate) (string, float64, error)
}

type Swipe interface {
	SaveSwipe(ctx context.Context, swipe *domain.Swipe) error
	GetSwipesCount(ctx context.Context) (int, error)
	GetSwipesByLobbyID(ctx context.Context, lobbyID string) ([]*domain.Swipe, error)
}

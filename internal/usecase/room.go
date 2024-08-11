package usecase

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"dishdash.ru/internal/domain"
	"dishdash.ru/pkg/filter"
)

type Match struct {
	ID    int
	Place *domain.Place
}

type Room struct {
	Lobby *domain.Lobby

	Users           map[string]*domain.User
	UsersPlace      map[string]*domain.Place
	UsersPlaceMutex sync.Mutex

	Places  []*domain.Place
	swipes  []*domain.Swipe
	matches []*Match

	usersMutex  sync.RWMutex
	placesMutex sync.RWMutex
	swipesMutex sync.RWMutex

	lobbyUseCase Lobby
	placeUseCase Place
}

func NewRoom(
	lobby *domain.Lobby,
	lobbyUseCase Lobby,
	placeUseCase Place,
) *Room {
	return &Room{
		Lobby:        lobby,
		Users:        make(map[string]*domain.User),
		Places:       make([]*domain.Place, 0),
		UsersPlace:   make(map[string]*domain.Place),
		swipes:       make([]*domain.Swipe, 0),
		usersMutex:   sync.RWMutex{},
		placesMutex:  sync.RWMutex{},
		swipesMutex:  sync.RWMutex{},
		lobbyUseCase: lobbyUseCase,
		placeUseCase: placeUseCase,
	}
}

func (r *Room) AddUser(user *domain.User) error {
	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()

	if _, has := r.Users[user.ID]; has {
		return fmt.Errorf("user %s already exists", user.ID)
	}
	r.Users[user.ID] = user
	return nil
}

func (r *Room) RemoveUser(id string) error {
	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()

	_, has := r.Users[id]
	if !has {
		return fmt.Errorf("user %s not found", id)
	}
	delete(r.Users, id)
	return nil
}

func (r *Room) UpdateLobby(ctx context.Context, input UpdateLobbyInput) error {
	lobby, err := r.lobbyUseCase.UpdateLobby(ctx, input)
	if err != nil {
		return err
	}
	r.Lobby = lobby
	return nil
}

func (r *Room) StartSwipes(ctx context.Context) error {
	r.swipesMutex.Lock()
	defer r.swipesMutex.Unlock()

	var err error
	r.Places, err = r.placeUseCase.GetPlacesForLobby(ctx, r.Lobby)
	if err != nil {
		return err
	}
	err = r.UpdateLobby(ctx, UpdateLobbyInput{
		ID: r.Lobby.ID,
		SaveLobbyInput: SaveLobbyInput{
			PriceAvg: r.Lobby.PriceAvg,
			Location: r.Lobby.Location,
			Tags: filter.Map(r.Lobby.Tags, func(t *domain.Tag) int64 {
				return t.ID
			}),
			Places: filter.Map(r.Places, func(p *domain.Place) int64 {
				return p.ID
			}),
		},
	})

	r.UsersPlaceMutex.Lock()
	defer r.UsersPlaceMutex.Unlock()
	for id := range r.Users {
		r.UsersPlace[id] = r.Places[0]
	}

	return err
}

func (r *Room) Swipe(userID string, placeID int64, t domain.SwipeType) (*Match, error) {
	r.swipesMutex.RLock()
	defer r.swipesMutex.RUnlock()

	r.swipes = append(r.swipes, &domain.Swipe{
		LobbyID: r.Lobby.ID,
		PlaceID: placeID,
		UserID:  userID,
		Type:    t,
	})

	matches := filter.Filter(r.swipes, func(swipe *domain.Swipe) bool {
		return swipe.PlaceID == placeID && swipe.Type == domain.LIKE
	})

	var match *Match

	if len(matches) > len(r.Users)/2 {
		match = &Match{Place: r.Places[slices.IndexFunc(r.Places, func(place *domain.Place) bool {
			return place.ID == placeID
		})]}
		match.ID = len(r.matches)
		r.matches = append(r.matches, match)
	}

	pIdx := slices.IndexFunc(r.Places, func(place *domain.Place) bool {
		return place.ID == placeID
	})

	r.UsersPlaceMutex.Lock()
	defer r.UsersPlaceMutex.Unlock()
	r.UsersPlace[userID] = r.Places[(pIdx+1)%len(r.Places)]

	return match, nil
}

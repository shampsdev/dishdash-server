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
	ID    string
	lobby *domain.Lobby

	users           map[string]*domain.User
	usersPlace      map[string]*domain.Place
	usersPlaceMutex sync.RWMutex

	places  []*domain.Place
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
		ID:           lobby.ID,
		lobby:        lobby,
		users:        make(map[string]*domain.User),
		places:       make([]*domain.Place, 0),
		usersPlace:   make(map[string]*domain.Place),
		swipes:       make([]*domain.Swipe, 0),
		usersMutex:   sync.RWMutex{},
		placesMutex:  sync.RWMutex{},
		swipesMutex:  sync.RWMutex{},
		lobbyUseCase: lobbyUseCase,
		placeUseCase: placeUseCase,
	}
}

func (r *Room) GetNextPlaceForUser(id string) *domain.Place {
	r.usersPlaceMutex.RLock()
	defer r.usersPlaceMutex.RUnlock()
	return r.usersPlace[id]
}

func (r *Room) AddUser(user *domain.User) error {
	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()

	if _, has := r.users[user.ID]; has {
		return fmt.Errorf("user %s already exists", user.ID)
	}
	r.users[user.ID] = user
	return nil
}

func (r *Room) RemoveUser(id string) error {
	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()

	_, has := r.users[id]
	if !has {
		return fmt.Errorf("user %s not found", id)
	}
	delete(r.users, id)
	return nil
}

func (r *Room) Empty() bool {
	return len(r.users) == 0
}

func (r *Room) UpdateLobby(ctx context.Context, priceAvg int, tagIDs []int64, placeIDs []int64) error {
	lobby, err := r.lobbyUseCase.UpdateLobby(ctx, UpdateLobbyInput{
		ID: r.lobby.ID,
		SaveLobbyInput: SaveLobbyInput{
			PriceAvg: priceAvg,
			Location: r.lobby.Location,
			Tags:     tagIDs,
			Places:   placeIDs,
		},
	})
	if err != nil {
		return err
	}
	r.lobby = lobby
	return nil
}

func (r *Room) StartSwipes(ctx context.Context) error {
	r.swipesMutex.Lock()
	defer r.swipesMutex.Unlock()

	var err error
	r.places, err = r.placeUseCase.GetPlacesForLobby(ctx, r.lobby)
	if err != nil {
		return err
	}
	err = r.UpdateLobby(ctx, r.lobby.PriceAvg,
		filter.Map(r.lobby.Tags, func(t *domain.Tag) int64 {
			return t.ID
		}),
		filter.Map(r.places, func(p *domain.Place) int64 {
			return p.ID
		}))

	r.usersPlaceMutex.Lock()
	defer r.usersPlaceMutex.Unlock()
	for id := range r.users {
		r.usersPlace[id] = r.places[0]
	}

	return err
}

func (r *Room) Swipe(userID string, placeID int64, t domain.SwipeType) (*Match, error) {
	r.swipesMutex.RLock()
	defer r.swipesMutex.RUnlock()

	r.swipes = append(r.swipes, &domain.Swipe{
		LobbyID: r.lobby.ID,
		PlaceID: placeID,
		UserID:  userID,
		Type:    t,
	})

	matches := filter.Filter(r.swipes, func(swipe *domain.Swipe) bool {
		return swipe.PlaceID == placeID && swipe.Type == domain.LIKE
	})

	var match *Match

	if len(matches) > len(r.users)/2 {
		match = &Match{Place: r.places[slices.IndexFunc(r.places, func(place *domain.Place) bool {
			return place.ID == placeID
		})]}
		match.ID = len(r.matches)
		r.matches = append(r.matches, match)
	}

	pIdx := slices.IndexFunc(r.places, func(place *domain.Place) bool {
		return place.ID == placeID
	})

	r.usersPlaceMutex.Lock()
	defer r.usersPlaceMutex.Unlock()
	r.usersPlace[userID] = r.places[(pIdx+1)%len(r.places)]

	return match, nil
}

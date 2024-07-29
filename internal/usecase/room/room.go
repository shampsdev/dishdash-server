package room

import (
	"context"
	"fmt"
	"slices"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
	"dishdash.ru/pkg/filter"
)

type Match struct {
	Place *domain.Place
}

type Room struct {
	Lobby  *domain.Lobby
	Users  map[string]*domain.User
	places []*domain.Place
	swipes []*domain.Swipe

	lobbyUseCase usecase.LobbyUseCase
	placeUseCase usecase.Place
}

func NewRoom(lobby *domain.Lobby) *Room {
	return &Room{lobby: lobby}
}

func (r *Room) AddUser(user *domain.User) error {
	if _, has := r.Users[user.ID]; has {
		return fmt.Errorf("user %s already exists", user.ID)
	}
	r.Users[user.ID] = user
	return nil
}

func (r *Room) RemoveUser(id string) error {
	_, has := r.Users[id]
	if !has {
		return fmt.Errorf("user %s not found", id)
	}
	delete(r.Users, id)
	return nil
}

func (r *Room) UpdateLobby(ctx context.Context, input usecase.UpdateLobbyInput) error {
	lobby, err := r.lobbyUseCase.UpdateLobby(ctx, input)
	if err != nil {
		return err
	}
	r.lobby = lobby
	return nil
}

func (r *Room) StartSwipes(ctx context.Context) error {
	var err error
	r.places, err = r.placeUseCase.GetAllPlaces(ctx)
	if err != nil {
		return err
	}
	err = r.UpdateLobby(ctx, usecase.UpdateLobbyInput{
		ID: r.lobby.ID,
		SaveLobbyInput: usecase.SaveLobbyInput{
			PriceAvg: r.lobby.PriceAvg,
			Location: r.lobby.Location,
			Tags: filter.Map(r.lobby.Tags, func(t *domain.Tag) int64 {
				return t.ID
			}),
			Places: filter.Map(r.places, func(p *domain.Place) int64 {
				return p.ID
			}),
		},
	})

	return err
}

func (r *Room) Swipe(userID string, placeID int64, t domain.SwipeType) (*Match, error) {
	r.swipes = append(r.swipes, &domain.Swipe{
		LobbyID: r.lobby.ID,
		PlaceID: placeID,
		UserID:  userID,
		Type:    t,
	})

	matches := filter.Filter(r.swipes, func(swipe *domain.Swipe) bool {
		return swipe.PlaceID == placeID && swipe.Type == domain.LIKE
	})

	match := new(Match)

	if len(matches) > len(r.Users)/2 {
		match = &Match{Place: r.places[slices.IndexFunc(r.places, func(place *domain.Place) bool {
			return place.ID == placeID
		})]}
	}

	return match, nil
}

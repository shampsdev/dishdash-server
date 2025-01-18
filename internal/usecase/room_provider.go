package usecase

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

type RoomRepo interface {
	GetRoom(ctx context.Context, id string) (*NRoom, error)
	GetActiveRoomCount() (int, error)
	DeleteRoom(ctx context.Context, id string) error
}

type InMemoryRoomRepo struct {
	lobbyUseCase     Lobby
	placesUseCase    Place
	swipeUseCase     Swipe
	userUseCase      User
	placeRecommender *PlaceRecommender

	roomsMutex sync.RWMutex
	rooms      map[string]*NRoom
}

func NewInMemoryRoomRepo(
	lobbyUseCase Lobby,
	placeUseCase Place,
	swipeUseCase Swipe,
	userUseCase User,
	placeRecomender *PlaceRecommender,
) *InMemoryRoomRepo {
	return &InMemoryRoomRepo{
		lobbyUseCase:     lobbyUseCase,
		placesUseCase:    placeUseCase,
		userUseCase:      userUseCase,
		placeRecommender: placeRecomender,
		swipeUseCase:     swipeUseCase,
		rooms:            make(map[string]*NRoom),
	}
}

func (r *InMemoryRoomRepo) GetRoom(ctx context.Context, id string) (*NRoom, error) {
	r.roomsMutex.Lock()
	defer r.roomsMutex.Unlock()
	room, ok := r.rooms[id]
	if !ok {
		lobby, err := r.lobbyUseCase.GetLobbyByID(ctx, id)
		if err != nil {
			return nil, err
		}

		log.Infof("Create room: %s", id)
		room, err := NewNRoom(
			lobby,
			r.lobbyUseCase,
			r.placesUseCase,
			r.swipeUseCase,
			r.userUseCase,
			r.placeRecommender,
		)
		if err != nil {
			return nil, err
		}
		r.rooms[id] = room
		return room, nil
	}

	return room, nil
}

func (r *InMemoryRoomRepo) GetActiveRoomCount() (int, error) {
	return len(r.rooms), nil
}

func (r *InMemoryRoomRepo) DeleteRoom(_ context.Context, id string) error {
	r.roomsMutex.Lock()
	defer r.roomsMutex.Unlock()
	delete(r.rooms, id)
	log.Infof("Deleted room: %s", id)
	return nil
}

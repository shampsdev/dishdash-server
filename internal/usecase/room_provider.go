package usecase

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type RoomRepo interface {
	GetRoom(ctx context.Context, id string) (*Room, error)
	GetActiveRoomCount() (int, error)
	DeleteRoom(ctx context.Context, id string) error
}

type InMemoryRoomRepo struct {
	lobbyUseCase     Lobby
	placesUseCase    Place
	swipeUseCase     Swipe
	userUseCase      User
	placeRecommender *PlaceRecommender

	roomsMutex  sync.RWMutex
	rooms       map[string]*Room
	activeRooms map[string]time.Time
}

func NewInMemoryRoomRepo(
	lobbyUseCase Lobby,
	placeUseCase Place,
	swipeUseCase Swipe,
	userUseCase User,
	placeRecomender *PlaceRecommender,
) *InMemoryRoomRepo {
	r := &InMemoryRoomRepo{
		lobbyUseCase:     lobbyUseCase,
		placesUseCase:    placeUseCase,
		userUseCase:      userUseCase,
		placeRecommender: placeRecomender,
		swipeUseCase:     swipeUseCase,
		rooms:            make(map[string]*Room),
		activeRooms:      make(map[string]time.Time),
	}

	r.goGC()
	return r
}

func (r *InMemoryRoomRepo) GetRoom(ctx context.Context, id string) (*Room, error) {
	r.roomsMutex.Lock()
	defer r.roomsMutex.Unlock()
	room, ok := r.rooms[id]
	if !ok {
		lobby, err := r.lobbyUseCase.GetLobbyByID(ctx, id)
		if err != nil {
			return nil, err
		}

		log.Infof("Create room: %s", id)
		room, err := NewRoom(
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
		r.activeRooms[id] = time.Now()
		return room, nil
	}

	r.activeRooms[id] = time.Now()
	return room, nil
}

func (r *InMemoryRoomRepo) goGC() {
	timer := time.NewTimer(time.Minute)
	go func() {
		for range timer.C {
			r.roomsMutex.Lock()
			for id, room := range r.rooms {
				if room.Active() {
					r.activeRooms[id] = time.Now()
				} else {
					wasActive := r.activeRooms[id]
					if time.Since(wasActive) > time.Minute*5 {
						delete(r.rooms, id)
						log.Infof("Deleted room: %s", id)
					}
				}
			}
			r.roomsMutex.Unlock()
		}
	}()
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

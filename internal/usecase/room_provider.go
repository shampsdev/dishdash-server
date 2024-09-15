package usecase

import (
	"context"
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"
)

type RoomRepo interface {
	GetRoom(ctx context.Context, id string) (*Room, error)
	DeleteRoom(ctx context.Context, id string) error
}

type InMemoryRoomRepo struct {
	lobbyUseCase  Lobby
	placesUseCase Place

	roomsMutex sync.RWMutex
	rooms      map[string]*Room
}

func NewInMemoryRoomRepo(lobbyUseCase Lobby, placeUseCase Place) *InMemoryRoomRepo {
	return &InMemoryRoomRepo{
		lobbyUseCase:  lobbyUseCase,
		placesUseCase: placeUseCase,
		rooms:         make(map[string]*Room),
	}
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
		room, err := NewRoom(lobby, r.lobbyUseCase, r.placesUseCase)
		if err != nil {
			return nil, err
		}
		r.rooms[id] = room
		return room, nil
	}

	if !room.InLobby() {
		return nil, errors.New("can't connect to started lobby")
	}

	return room, nil
}

func (r *InMemoryRoomRepo) DeleteRoom(_ context.Context, id string) error {
	r.roomsMutex.Lock()
	defer r.roomsMutex.Unlock()
	delete(r.rooms, id)
	log.Infof("Deleted room: %s", id)
	return nil
}

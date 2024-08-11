package usecase

import (
	"context"
	"log"
	"sync"
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

		log.Printf("create room: %s", id)
		room := NewRoom(lobby, r.lobbyUseCase, r.placesUseCase)
		r.rooms[id] = room
		return room, nil
	}
	return room, nil
}

func (r *InMemoryRoomRepo) DeleteRoom(_ context.Context, id string) error {
	r.roomsMutex.Lock()
	defer r.roomsMutex.Unlock()
	delete(r.rooms, id)
	log.Printf("deleted room: %s", id)
	return nil
}
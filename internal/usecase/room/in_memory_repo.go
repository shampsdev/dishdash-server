package room

import (
	"context"
	"log"

	"dishdash.ru/internal/usecase"
)

type InMemoryRepo struct {
	lobbyUseCase usecase.Lobby
	rooms        map[string]*Room
}

func NewInMemoryRepo(lobbyUseCase usecase.Lobby) *InMemoryRepo {
	return &InMemoryRepo{lobbyUseCase: lobbyUseCase}
}

func (r InMemoryRepo) GetRoom(ctx context.Context, id string) (*Room, error) {
	room, ok := r.rooms[id]
	if !ok {
		lobby, err := r.lobbyUseCase.GetLobbyByID(ctx, id)

		if err != nil {
			return nil, err
		}

		log.Printf("create room: %s", id)
		return NewRoom(lobby), err
	}
	return room, nil
}

func (r InMemoryRepo) DeleteRoom(_ context.Context, id string) error {
	delete(r.rooms, id)
	log.Printf("deleted room: %s", id)
	return nil
}

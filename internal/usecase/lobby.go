package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Lobby struct {
	lobbyRepo LobbyRepository
}

func NewLobby(lobbyRepo LobbyRepository) *Lobby {
	return &Lobby{lobbyRepo: lobbyRepo}
}

func (lb *Lobby) GetLobbyByID(ctx context.Context, id int64) (*domain.Lobby, error) {
	return lb.lobbyRepo.GetLobbyByID(ctx, id)
}

func (lb *Lobby) SaveLobby(ctx context.Context, lobby *domain.Lobby) error {
	return lb.lobbyRepo.SaveLobby(ctx, lobby)
}

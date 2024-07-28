package usecase

import (
	"context"
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"
	"errors"
)

type LobbyUseCase struct {
	lRepo repo.Lobby
	uRepo repo.User
	tRepo repo.Tag
	pRepo repo.Place
	sRepo repo.Swipe
}

func NewLobbyUseCase(
	lRepo repo.Lobby,
	uRepo repo.User,
	tRepo repo.Tag,
	pRepo repo.Place,
	sRepo repo.Swipe,
) *LobbyUseCase {
	return &LobbyUseCase{
		lRepo: lRepo,
		uRepo: uRepo,
		tRepo: tRepo,
		pRepo: pRepo,
		sRepo: sRepo,
	}
}

func (l LobbyUseCase) SaveLobby(ctx context.Context, lobbyInput SaveLobbyInput) (*domain.Lobby, error) {
	lobby := &domain.Lobby{
		State:    domain.ACTIVE,
		PriceAvg: lobbyInput.PriceAvg,
		Location: lobbyInput.Location,
	}
	id, err := l.lRepo.SaveLobby(ctx, lobby)
	if err != nil {
		return nil, err
	}
	lobby.ID = id

	err = l.tRepo.AttachTagsToLobby(ctx, lobbyInput.Tags, id)
	if err != nil {
		return nil, err
	}

	return l.GetLobbyByID(ctx, id)
}

func (l LobbyUseCase) UpdateLobby(ctx context.Context, lobbyInput UpdateLobbyInput) (*domain.Lobby, error) {
	lobby := &domain.Lobby{
		ID:       lobbyInput.ID,
		State:    domain.ACTIVE,
		PriceAvg: lobbyInput.PriceAvg,
		Location: lobbyInput.Location,
	}
	err := l.lRepo.UpdateLobby(ctx, lobby)
	if err != nil {
		return nil, err
	}

	err = l.tRepo.DetachTagsFromLobby(ctx, lobby.ID)
	if err != nil {
		return nil, err
	}

	err = l.tRepo.AttachTagsToLobby(ctx, lobbyInput.Tags, lobby.ID)
	if err != nil {
		return nil, err
	}

	return l.GetLobbyByID(ctx, lobby.ID)
}

func (l LobbyUseCase) DeleteLobbyByID(ctx context.Context, id string) error {
	return l.lRepo.DeleteLobbyByID(ctx, id)
}

func (l LobbyUseCase) GetLobbyByID(ctx context.Context, id string) (*domain.Lobby, error) {
	lobby, err := l.lRepo.GetLobbyByID(ctx, id)
	if err != nil {
		return nil, err
	}
	lobby.ID = id

	lobby.Tags, err = l.tRepo.GetTagsByLobbyID(ctx, id)
	if err != nil {
		return nil, err
	}

	lobby.Swipes, err = l.sRepo.GetSwipesByLobbyID(ctx, id)
	if err != nil {
		return nil, err
	}

	lobby.Users, err = l.uRepo.GetUsersByLobbyID(ctx, id)
	if err != nil {
		return nil, err
	}

	lobby.Places, err = l.pRepo.GetPlacesByLobbyID(ctx, id)
	if err != nil {
		return nil, err
	}

	return lobby, nil
}

func (l LobbyUseCase) NearestActiveLobby(ctx context.Context, loc domain.Coordinate) (*domain.Lobby, float64, error) {
	id, dist, err := l.lRepo.NearestActiveLobbyID(ctx, loc)
	if err != nil {
		return nil, 0, err
	}
	lobby, err := l.GetLobbyByID(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	return lobby, dist, nil
}

func (l LobbyUseCase) FindLobby(ctx context.Context, input FindLobbyInput) (*domain.Lobby, error) {
	lobby, dist, err := l.NearestActiveLobby(ctx, input.Location)
	if err != nil && !errors.Is(err, repo.ErrLobbyNotFound) {
		return nil, err
	}
	if dist > input.Dist || errors.Is(err, repo.ErrLobbyNotFound) {
		lobby, err = l.SaveLobby(ctx, SaveLobbyInput{
			Location: input.Location,
			PriceAvg: 500,
		})
		if err != nil {
			return nil, err
		}
	}
	return l.GetLobbyByID(ctx, lobby.ID)
}

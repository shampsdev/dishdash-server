package usecase

import (
	"context"
	"errors"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"
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

func (l LobbyUseCase) SaveLobby(ctx context.Context, lobbyInput SaveLobbyInput) (*LobbyOutput, error) {
	lobby := &domain.Lobby{
		State:    domain.InLobby,
		PriceAvg: lobbyInput.PriceAvg,
		Location: lobbyInput.Location,
	}
	id, err := l.lRepo.SaveLobby(ctx, lobby)
	if err != nil {
		return nil, err
	}
	lobby.ID = id

	return l.GetOutputLobbyByID(ctx, id)
}

func (l LobbyUseCase) SetLobbySettings(ctx context.Context, lobbyInput UpdateLobbySettingsInput) (*domain.Lobby, error) {
	lobby := &domain.Lobby{
		ID:       lobbyInput.ID,
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

	err = l.pRepo.DetachPlacesFromLobby(ctx, lobby.ID)
	if err != nil {
		return nil, err
	}

	err = l.pRepo.AttachOrderedPlacesToLobby(ctx, lobbyInput.Places, lobby.ID)
	if err != nil {
		return nil, err
	}

	return l.GetLobbyByID(ctx, lobby.ID)
}

func (l LobbyUseCase) SetLobbyState(ctx context.Context, lobbyID string, state domain.LobbyState) error {
	return l.lRepo.SetLobbyState(ctx, lobbyID, state)
}

func (l LobbyUseCase) SetLobbyUsers(ctx context.Context, lobbyID string, userIDs []string) ([]*domain.User, error) {
	err := l.uRepo.DetachUsersFromLobby(ctx, lobbyID)
	if err != nil {
		return nil, err
	}
	err = l.uRepo.AttachUsersToLobby(ctx, userIDs, lobbyID)
	if err != nil {
		return nil, err
	}

	lobby, err := l.lRepo.GetLobbyByID(ctx, lobbyID)
	if err != nil {
		return nil, err
	}

	return lobby.Users, nil
}

func (l LobbyUseCase) DeleteLobbyByID(ctx context.Context, id string) error {
	err := l.tRepo.DetachTagsFromLobby(ctx, id)
	if err != nil {
		return err
	}

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

	lobby.Places, err = l.pRepo.GetOrderedPlacesByLobbyID(ctx, id)
	if err != nil {
		return nil, err
	}

	return lobby, nil
}

func (l LobbyUseCase) GetOutputLobbyByID(ctx context.Context, id string) (*LobbyOutput, error) {
	lobby, err := l.lRepo.GetLobbyByID(ctx, id)
	if err != nil {
		return nil, err
	}
	lobby.ID = id

	lobby.Tags, err = l.tRepo.GetTagsByLobbyID(ctx, id)
	if err != nil {
		return nil, err
	}

	lobby.Users, err = l.uRepo.GetUsersByLobbyID(ctx, id)
	if err != nil {
		return nil, err
	}

	lobbyOutput := &LobbyOutput{
		ID:        lobby.ID,
		State:     lobby.State,
		PriceAvg:  lobby.PriceAvg,
		Location:  lobby.Location,
		CreatedAt: lobby.CreatedAt,
		Tags:      lobby.Tags,
		Users:     lobby.Users,
	}

	return lobbyOutput, nil
}

func (l LobbyUseCase) NearestActiveLobby(ctx context.Context, loc domain.Coordinate) (*LobbyOutput, float64, error) {
	id, dist, err := l.lRepo.NearestActiveLobbyID(ctx, loc)
	if err != nil {
		return nil, 0, err
	}
	lobby, err := l.GetOutputLobbyByID(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	return lobby, dist, nil
}

func (l LobbyUseCase) FindLobby(ctx context.Context, input FindLobbyInput) (*LobbyOutput, error) {
	lobby, dist, err := l.NearestActiveLobby(ctx, input.Location)
	if err != nil && !errors.Is(err, repo.ErrLobbyNotFound) {
		return nil, err
	}
	if dist > input.Dist || errors.Is(err, repo.ErrLobbyNotFound) {
		lobby, err = l.SaveLobby(ctx, SaveLobbyInput{
			Location: input.Location,
			PriceAvg: 1200,
		})
		if err != nil {
			return nil, err
		}
	}
	return l.GetOutputLobbyByID(ctx, lobby.ID)
}

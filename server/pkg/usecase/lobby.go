package usecase

import (
	"context"
	"fmt"

	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/repo"
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

func (l LobbyUseCase) CreateLobby(ctx context.Context, settings domain.LobbySettings) (*domain.Lobby, error) {
	lobby := &domain.Lobby{
		State:    domain.InLobby,
		Type:     settings.Type,
		Settings: settings,
	}
	var err error
	settings, err = l.validateLobbySettings(settings)
	if err != nil {
		return nil, fmt.Errorf("invalid lobby settings: %w", err)
	}

	id, err := l.lRepo.SaveLobby(ctx, lobby)
	if err != nil {
		return nil, fmt.Errorf("failed to save lobby: %w", err)
	}
	lobby.ID = id

	return l.GetLobbyByID(ctx, id)
}

func (l LobbyUseCase) validateLobbySettings(settings domain.LobbySettings) (domain.LobbySettings, error) {
	switch settings.Type {
	case domain.ClassicPlacesLobbyType:
		if settings.ClassicPlaces == nil {
			return domain.LobbySettings{}, fmt.Errorf("classic places settings are required")
		}

		if settings.ClassicPlaces.Recommendation == nil {
			settings.ClassicPlaces.Recommendation = defaultRecommendationOpts()
		}

		return settings, nil
	default:
		return domain.LobbySettings{}, fmt.Errorf("unsupported lobby type: %s", settings.Type)
	}
}

func (l LobbyUseCase) SetLobbySettings(ctx context.Context, lobbyID string, settings domain.LobbySettings) error {
	return l.lRepo.SetLobbySettings(ctx, lobbyID, settings)
}

func (l LobbyUseCase) AttachOrderedPlacesToLobby(ctx context.Context, placeIDs []int64, lobbyID string) error {
	return l.pRepo.AttachOrderedPlacesToLobby(ctx, placeIDs, lobbyID)
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

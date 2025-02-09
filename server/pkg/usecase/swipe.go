package usecase

import (
	"context"

	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/repo"
)

type SwipeUseCase struct {
	sRepo repo.Swipe
}

func NewSwipeUseCase(sRepo repo.Swipe) *SwipeUseCase {
	return &SwipeUseCase{sRepo: sRepo}
}

func (s *SwipeUseCase) SaveSwipe(ctx context.Context, swipe *domain.Swipe) error {
	return s.sRepo.SaveSwipe(ctx, swipe)
}

func (s *SwipeUseCase) GetCount(ctx context.Context) (int, error) {
	return s.sRepo.GetSwipesCount(ctx)
}

func (s *SwipeUseCase) GetSwipesByLobbyID(ctx context.Context, lobbyID string) ([]*domain.Swipe, error) {
	return s.sRepo.GetSwipesByLobbyID(ctx, lobbyID)
}

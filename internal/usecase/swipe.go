package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"
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
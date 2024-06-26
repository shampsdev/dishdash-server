package usecase

import (
	"context"
	"log"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"
)

type SwipeUseCase struct {
	swipeRepo repo.Swipe
}

func NewSwipeUseCase(swipeRepo repo.Swipe) *SwipeUseCase {
	log.Println("Usecase created")
	return &SwipeUseCase{swipeRepo: swipeRepo}
}

func (s *SwipeUseCase) CreateSwipe(ctx context.Context, swipe *domain.Swipe) error {
	return s.swipeRepo.CreateSwipe(ctx, swipe)
}

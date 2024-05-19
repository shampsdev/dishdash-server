package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
)

type Swipe struct {
	swipeRepo SwipeRepository
}

func NewSwipe(swipeRepo SwipeRepository) *Swipe {
	return &Swipe{swipeRepo: swipeRepo}
}

func (s *Swipe) SaveSwipe(ctx context.Context, swipe *domain.Swipe) error {
	return s.SaveSwipe(ctx, swipe)
}

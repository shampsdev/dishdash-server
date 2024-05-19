package swipe

import (
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/dto"
)

type swipe struct {
	T    dto.SwipeType
	Card *domain.Card
}

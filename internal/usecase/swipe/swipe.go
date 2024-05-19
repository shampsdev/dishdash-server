package swipe

import (
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/dto"
)

type Swipe struct {
	T    dto.SwipeType
	Card *domain.Card
}

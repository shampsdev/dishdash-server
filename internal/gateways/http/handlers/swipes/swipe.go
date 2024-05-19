package swipes

import "dishdash.ru/internal/domain"

type swipe struct {
	T    swipeType
	Card *domain.Card
}

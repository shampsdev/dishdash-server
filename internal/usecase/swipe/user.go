package swipe

import (
	"context"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/dto"
	"dishdash.ru/internal/usecase"
)

type User struct {
	swipeUseCase *usecase.Swipe
	id           string
	Lobby        *Lobby
	took         int
}

func NewUser(id string, swipeUseCase *usecase.Swipe) *User {
	return &User{id: id, swipeUseCase: swipeUseCase}
}

func (u *User) Card() *domain.Card {
	return u.Lobby.takeCard(u.took)
}

// Swipe returns matched card if was match
func (u *User) Swipe(swipeType dto.SwipeType) *domain.Card {
	_ = u.swipeUseCase.SaveSwipe(context.Background(), &domain.Swipe{
		LobbyID:   u.Lobby.id,
		CardID:    u.Card().ID,
		UserID:    u.id,
		SwipeType: swipeType,
	})
	if swipeType == dto.LIKE {
		card := u.Lobby.like(u.Card())
		u.took++
		return card
	}
	u.took++
	return nil
}

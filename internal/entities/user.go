package entities

import (
	"context"
	"log"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
)

type User struct {
	swipeUseCase usecase.Swipe
	ID           string
	Lobby        *Lobby
	took         int
}

func NewUser(user domain.User, swipeUseCase usecase.Swipe) *User {
	return &User{ID: user.ID, swipeUseCase: swipeUseCase}
}

func (u *User) Card() *domain.Card {
	return u.Lobby.takeCard(u.took)
}

func (u *User) Swipe(swipeType domain.SwipeType) *domain.Card {
	if u.took >= len(u.Lobby.cards) {
		u.took++
		return nil
	}

	log.Println(u.swipeUseCase)

	err := u.swipeUseCase.CreateSwipe(context.Background(), &domain.Swipe{
		LobbyID: u.Lobby.ID,
		CardID:  u.Card().ID,
		UserID:  u.ID,
		Type:    swipeType,
	})
	if err != nil {
		log.Println("Swipe wasn't able to be created")
	}

	if swipeType == domain.LIKE {
		card := u.Lobby.like(u.Card())
		u.took++
		return card
	}
	u.took++

	return nil
}

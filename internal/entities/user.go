package entities

import (
	"context"
	"log"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
)

type User struct {
	swipeUseCase usecase.Swipe
	id           string
	Lobby        *Lobby
	took         int
}

func NewUser(user domain.User, swipeUseCase usecase.Swipe) *User {
	log.Println("In the user its this", swipeUseCase)
	return &User{id: user.ID, swipeUseCase: swipeUseCase}
}

func (u *User) Card() *domain.Card {
	return u.Lobby.takeCard(u.took)
}

// Swipe returns matched card if was match
func (u *User) Swipe(swipeType domain.SwipeType) *domain.Card {
	// TODO better cards end logic
	if u.took >= len(u.Lobby.cards) {
		u.took++
		return nil
	}

	log.Println(u.swipeUseCase)

	u.swipeUseCase.CreateSwipe(context.Background(), &domain.Swipe{
		LobbyID: u.Lobby.Id,
		CardID:  u.Card().ID,
		UserID:  u.id,
		Type:    swipeType,
	})

	if swipeType == domain.LIKE {
		card := u.Lobby.like(u.Card())
		u.took++
		return card
	}
	u.took++

	log.Println("Problem is here lol 4")

	return nil
}

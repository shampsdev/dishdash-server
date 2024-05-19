package swipes

import (
	"context"
	"sync"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
)

type lobby struct {
	*domain.Lobby
	cards []*domain.Card
	users map[*user]bool

	lock sync.RWMutex
}

var lobbies = make(map[int64]*lobby)

func findLobby(domainLobby *domain.Lobby, cardUseCase *usecase.Card) (*lobby, error) {
	lb, has := lobbies[domainLobby.ID]
	if has {
		return lb, nil
	}

	lb = &lobby{
		Lobby: domainLobby,
		users: make(map[*user]bool),
	}
	var err error
	lb.cards, err = cardUseCase.GetCards(context.Background())
	if err != nil {
		return nil, err
	}

	lb.lock.Lock()
	defer lb.lock.Unlock()
	lobbies[domainLobby.ID] = lb

	return lb, nil
}

func (lb *lobby) registerUser(u *user) {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	lb.users[u] = true
}

func (lb *lobby) unregisterUser(u *user) {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	delete(lb.users, u)
}

func (lb *lobby) takeCard(n int) *domain.Card {
	lb.lock.RLock()
	defer lb.lock.RUnlock()

	return lb.cards[n%len(lb.cards)]
}

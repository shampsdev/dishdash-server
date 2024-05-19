package swipe

import (
	"context"
	"log"
	"slices"
	"sync"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
)

type Lobby struct {
	Id    int64
	cards []*domain.Card
	likes map[*domain.Card]int
	users map[*User]bool

	lock sync.RWMutex
}

var lobbies = make(map[int64]*Lobby)

func FindLobby(lobbyDomain *domain.Lobby, cardUseCase *usecase.Card) (*Lobby, error) {
	lobby, has := lobbies[lobbyDomain.ID]
	log.Printf("find lobby: %d", lobbyDomain.ID)
	if has {
		return lobby, nil
	}

	cards, err := cardUseCase.GetCards(context.Background())

	slices.SortFunc(cards, func(a, b *domain.Card) int {
		d1 := lobbyDomain.Location.GreatCircleDistance(a.Location)
		d2 := lobbyDomain.Location.GreatCircleDistance(b.Location)
		return int(d1 - d2)
	})

	if len(cards) == 0 {
		panic("no cards in database")
	}
	if err != nil {
		return nil, err
	}
	lobby = &Lobby{
		Id:    lobbyDomain.ID,
		cards: cards,
		likes: make(map[*domain.Card]int),
		users: make(map[*User]bool),
		lock:  sync.RWMutex{},
	}
	log.Printf("create lobby: %d", lobbyDomain.ID)
	return lobby, nil
}

func (lb *Lobby) Register(u *User) {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	u.Lobby = lb
	lb.users[u] = true
	log.Printf("register user in lobby: %d", lb.Id)
}

func (lb *Lobby) Unregister(u *User) {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	delete(lb.users, u)
	log.Printf("unregister user in lobby: %d", lb.Id)
	if len(lb.users) == 0 {
		delete(lobbies, lb.Id)
		log.Printf("delete lobby: %d", lb.Id)
	}
}

func (lb *Lobby) takeCard(n int) *domain.Card {
	return lb.cards[n%len(lb.cards)]
}

// like returns card if match
func (lb *Lobby) like(card *domain.Card) *domain.Card {
	_, has := lb.likes[card]
	if !has {
		lb.likes[card] = 0
	}
	lb.likes[card]++
	if lb.likes[card] >= len(lb.users)/2 {
		return card
	}
	return nil
}

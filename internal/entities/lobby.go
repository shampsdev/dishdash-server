package entities

import (
	"context"
	"log"
	"sync"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
)

type Lobby struct {
	Id       string
	cards    []*domain.Card
	likes    map[*domain.Card]int
	users    map[string]*User
	settings domain.LobbySettings

	lock sync.RWMutex
}

var (
	lobbies     = make(map[string]*Lobby)
	lobbiesLock = sync.Mutex{}
)

func FindLobby(lobbyDomain *domain.Lobby, cardUseCase usecase.Card) (*Lobby, error) {
	lobbiesLock.Lock()
	defer lobbiesLock.Unlock()

	lobby, has := lobbies[lobbyDomain.ID]
	if has {
		log.Printf("find lobby: %s", lobbyDomain.ID)
		return lobby, nil
	}

	cards, err := cardUseCase.GetAllCards(context.Background())

	// slices.SortFunc(cards, func(a, b *domain.Card) int {
	// 	d1 := lobbyDomain.Location.GreatCircleDistance(a.Location)
	// 	d2 := lobbyDomain.Location.GreatCircleDistance(b.Location)
	// 	return int(d1 - d2)
	// })

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
		users: make(map[string]*User),
		lock:  sync.RWMutex{},
		settings: domain.LobbySettings{
			PriceMin:    0,
			PriceMax:    0,
			MaxDistance: 0,
			Tags:        []domain.Tag{},
		},
	}
	lobbies[lobbyDomain.ID] = lobby
	log.Printf("create lobby: %s", lobbyDomain.ID)
	return lobby, nil
}

func (lb *Lobby) GetUsers() []*User {
	lb.lock.RLock() // Use RLock for concurrent read access
	defer lb.lock.RUnlock()

	var usersSlice []*User
	for _, user := range lb.users {
		usersSlice = append(usersSlice, user)
	}
	return usersSlice
}

func (lb *Lobby) UpdateSettings(settings domain.LobbySettings) {
	lb.settings = settings
}

func (lb *Lobby) Register(connectionId string, user *User) *User {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	user.Lobby = lb

	lb.users[connectionId] = user
	log.Printf("register user in lobby: %s", lb.Id)

	return user
}

func (lb *Lobby) Unregister(connectionId string) {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	delete(lb.users, connectionId)
	log.Printf("unregister user in lobby: %s", lb.Id)
	if len(lb.users) == 0 {
		delete(lobbies, lb.Id)
		log.Printf("delete lobby: %s", lb.Id)
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
	if lb.likes[card] > len(lb.users)/2 {
		return card
	}
	return nil
}

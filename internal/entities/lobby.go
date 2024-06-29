package entities

import (
	"context"
	"log"
	"slices"
	"sync"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
	"dishdash.ru/pkg/location"
)

type Lobby struct {
	ID       string
	Location domain.Coordinate
	cards    []*domain.Card
	likes    map[*domain.Card]int
	users    map[string]*User
	Settings domain.LobbySettings
	votes    map[int64]*Vote

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

	// Лучше это будет вынести в юз кейс или чет такое,
	// я кстати вообще не понимаю смысл их, я бы их сервисами сделал
	// ладно potato pothato
	// mike, 27/06/24

	cards, err := cardUseCase.GetAllCards(context.Background())
	slices.SortFunc(cards, func(a, b *domain.Card) int {
		d1 := location.GetDistance(lobbyDomain.Location, a.Location)
		d2 := location.GetDistance(lobbyDomain.Location, b.Location)
		return int(d1 - d2)
	})

	if len(cards) == 0 {
		panic("no cards in database")
	}
	if err != nil {
		return nil, err
	}
	lobby = &Lobby{
		ID:       lobbyDomain.ID,
		Location: lobbyDomain.Location,
		cards:    cards,
		likes:    make(map[*domain.Card]int),
		users:    make(map[string]*User),
		lock:     sync.RWMutex{},
		votes:    make(map[int64]*Vote),
		Settings: domain.LobbySettings{
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

func (lb *Lobby) SetCards(cards []*domain.Card) {
	lb.cards = cards
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
	lb.Settings = settings
}

func (lb *Lobby) Register(connectionID string, user *User) *User {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	user.Lobby = lb

	lb.users[connectionID] = user
	log.Printf("register user in lobby: %s", lb.ID)

	return user
}

func (lb *Lobby) Unregister(connectionID string) {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	delete(lb.users, connectionID)
	log.Printf("unregister user in lobby: %s", lb.ID)
	if len(lb.users) == 0 {
		delete(lobbies, lb.ID)
		log.Printf("delete lobby: %s", lb.ID)
	}
}

func (lb *Lobby) RegisterVote(vote *Vote, matchID int64) {
	log.Println("registered a vote", matchID)
	lb.votes[matchID] = vote
}

func (lb *Lobby) GetVoteByID(id int64) *Vote {
	log.Println("getting the vote", id)
	return lb.votes[id]
}

func (lb *Lobby) GetResults() []domain.Card {
	if len(lb.likes) == 0 {
		return nil
	}

	// Find the maximum number of likes
	maxLikes := 0
	for _, likes := range lb.likes {
		if likes > maxLikes {
			maxLikes = likes
		}
	}

	// Collect all cards with the maximum number of likes
	var result []domain.Card
	for card, likes := range lb.likes {
		if likes == maxLikes {
			result = append(result, *card)
		}
	}

	return result
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
	if lb.likes[card] > len(lb.users)/2 {
		return nil
	}
	lb.likes[card]++
	if lb.likes[card] > len(lb.users)/2 {
		return card
	}
	return nil
}

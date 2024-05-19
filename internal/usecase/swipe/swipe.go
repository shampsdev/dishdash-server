package swipe

import (
	"sync"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/dto"
)

type Swipe struct {
	T    dto.SwipeType
	Card *domain.Card
}

type Lobby struct {
	cards []*domain.Card
	users map[*User]bool

	lock sync.RWMutex
}

type User struct {
	id    string
	lobby *Lobby
}

var lobbies = make(map[int64]*Lobby)

func (lb *Lobby) FindLobby(id int64) *Lobby {
	lobby, has := lobbies[id]
	if has {
		return lobby
	}

	lobby = &Lobby{}
	return lobby
}

func (lb *Lobby) Register(u *User) {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	u.lobby = lb
	lb.users[u] = true
}

func (lb *Lobby) Unregister(u *User) {
	lb.lock.Lock()
	defer lb.lock.Unlock()
	delete(lb.users, u)
}

// Swipe returns matched card if was match
func (u *User) Swipe() *domain.Card {
	return nil
}

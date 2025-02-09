package domain

import (
	"time"

	"dishdash.ru/pkg/algo"
)

type Lobby struct {
	ID        string
	State     LobbyState
	PriceAvg  int
	Location  Coordinate
	CreatedAt time.Time

	Tags   []*Tag
	Swipes []*Swipe
	Users  []*User
	Places []*Place
}

func (l *Lobby) TagNames() []string {
	return algo.Map(l.Tags, func(t *Tag) string {
		return t.Name
	})
}

type LobbyState string

var (
	InLobby  LobbyState = "lobby"
	Swiping  LobbyState = "swiping"
	Finished LobbyState = "finished"
)

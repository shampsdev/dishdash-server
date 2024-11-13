package domain

import (
	"time"

	"dishdash.ru/pkg/filter"
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
	return filter.Map(l.Tags, func(t *Tag) string {
		return t.Name
	})
}

type LobbyState string

var (
	InLobby  LobbyState = "lobby"
	Swiping  LobbyState = "swiping"
	Voting   LobbyState = "voting"
	Finished LobbyState = "finished"
)

package domain

import (
	"time"

	"dishdash.ru/pkg/filter"
)

type Lobby struct {
	ID        string     `json:"id"`
	State     LobbyState `json:"state"`
	PriceAvg  int        `json:"priceAvg"`
	Location  Coordinate `json:"location"`
	CreatedAt time.Time  `json:"createdAt"`

	Tags   []*Tag   `json:"tags"`
	Swipes []*Swipe `json:"swipes"`
	Users  []*User  `json:"users"`
	Places []*Place `json:"places"`
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

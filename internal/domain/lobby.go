package domain

import (
	"time"
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
	var names []string
	for _, tag := range l.Tags {
		if tag != nil {
			names = append(names, tag.Name)
		}
	}
	return names
}

type LobbyState string

var (
	ACTIVE   LobbyState = "active"
	INACTIVE LobbyState = "inactive"
)

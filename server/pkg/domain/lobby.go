package domain

import (
	"time"
)

type Lobby struct {
	ID        string     `json:"id"`
	State     LobbyState `json:"state"`
	CreatedAt time.Time  `json:"createdAt"`

	Type     LobbyType     `json:"type"`
	Settings LobbySettings `json:"settings"`

	Swipes []*Swipe `json:"swipes"`
	Users  []*User  `json:"users"`
	Places []*Place `json:"places"`
}

type LobbyState string

var (
	InLobby LobbyState = "lobby"
	Swiping LobbyState = "swiping"
)

type LobbyType string

var (
	ClassicPlacesLobbyType    LobbyType = "classicPlaces"
	CollectionPlacesLobbyType LobbyType = "collectionPlaces"
)

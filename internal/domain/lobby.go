package domain

import (
	"time"
)

type Lobby struct {
	ID        string
	State     LobbyState
	CreatedAt time.Time

	Type     LobbyType
	Settings LobbySettings

	Swipes []*Swipe
	Users  []*User
	Places []*Place
}

type LobbyState string

var (
	InLobby  LobbyState = "lobby"
	Swiping  LobbyState = "swiping"
	Finished LobbyState = "finished"
)

type LobbyType string

var (
	ClassicPlacesLobbyType    LobbyType = "classicPlaces"
	CollectionPlacesLobbyType LobbyType = "collectionPlaces"
)

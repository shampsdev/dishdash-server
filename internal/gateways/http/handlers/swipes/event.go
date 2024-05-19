package swipes

import "dishdash.ru/internal/dto"

type EventType string

const (
	eventJoinLobby = "joinLobby"
	eventCard      = "card"
	eventSwipe     = "swipe"
	eventMatch     = "match"
)

type swipeType string

const (
	like    swipeType = "like"
	dislike swipeType = "dislike"
)

type swipeEvent struct {
	SwipeType swipeType `json:"swipeType"`
}

type matchEvent struct {
	Card dto.Card `json:"card"`
}

type cardEvent struct {
	Card dto.Card `json:"card"`
}

type joinLobbyEvent struct {
	LobbyID int64 `json:"lobbyID"`
}

package swipe

import "dishdash.ru/internal/dto"

type EventType string

const (
	eventJoinLobby = "joinLobby"
	eventCard      = "card"
	eventSwipe     = "swipe"
	eventMatch     = "match"
)

type swipeEvent struct {
	SwipeType dto.SwipeType `json:"swipeType"`
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

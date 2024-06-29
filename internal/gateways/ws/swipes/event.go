package swipes

import (
	"dishdash.ru/internal/domain"
)

const (
	eventJoinLobby      = "joinLobby"
	eventSettingsUpdate = "settingsUpdate"
	eventCard           = "card"
	eventSwipe          = "swipe"
	eventMatch          = "match"
	eventRelaseMatch    = "releaseMatch"
	eventUserJoined     = "userJoined"
	eventUserLeft       = "userLeft"
	eventStartSwipes    = "startSwipes"
	eventVote           = "vote"
	eventVoted          = "voted"
	eventFinish         = "finish"
	eventFinalVote      = "finalVote"
)

type swipeEvent struct {
	SwipeType domain.SwipeType `json:"swipeType"`
}

type matchEvent struct {
	ID   int64       `json:"id"`
	Card domain.Card `json:"card"`
}

type cardEvent struct {
	Card domain.Card `json:"card"`
}

type settingsUpdateEvent struct {
	PriceMin    int          `json:"priceMin"`
	PriceMax    int          `json:"priceMax"`
	MaxDistance float64      `json:"maxDistance"`
	Tags        []domain.Tag `json:"tags"`
}

type joinLobbyEvent struct {
	LobbyID string `json:"lobbyId"`
	UserID  string `json:"userId"`
}

type userJoinEvent struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type userLeftEvent struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type voteEvent struct {
	VoteID     int64 `json:"id"`
	VoteOption int64 `json:"option"`
}

type User struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type votedEvent struct {
	User       User  `json:"user"`
	VoteID     int64 `json:"id"`
	VoteOption int64 `json:"option"`
}

type finalVoteEvent struct {
	VoteID  int64         `json:"id"`
	Options []domain.Card `json:"options"`
}

type finishEvent struct {
	Result domain.Card `json:"result"`
}

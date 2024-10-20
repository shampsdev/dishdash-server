package event

import (
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
)

const (
	JoinLobby      = "joinLobby"
	UserJoined     = "userJoined"
	UserLeft       = "userLeft"
	SettingsUpdate = "settingsUpdate"
	StartSwipes    = "startSwipes"
	Place          = "card"
	Swipe          = "swipe"
	Match          = "match"
	Vote           = "vote"
	Voted          = "voted"
	ReleaseMatch   = "releaseMatch"
	Finish         = "finish"

	Error = "error"
)

type JoinLobbyEvent struct {
	LobbyID string `json:"lobbyId"`
	UserID  string `json:"userId"`
}

type UserJoinedEvent struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type UserLeftEvent struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type SettingsUpdateEvent struct {
	Location    domain.Coordinate `json:"location"`
	UserID      string            `json:"userId"`
	PriceMin    int               `json:"priceMin"`
	PriceMax    int               `json:"priceMax"`
	MaxDistance int               `json:"maxDistance"`
	Tags        []int64           `json:"tags"`
}

type PlaceEvent struct {
	ID   int64         `json:"id"`
	Card *domain.Place `json:"card"`
}

type SwipeEvent struct {
	SwipeType domain.SwipeType `json:"swipeType"`
}

type MatchEvent struct {
	ID   int           `json:"id"`
	Card *domain.Place `json:"card"`
}

type VoteEvent struct {
	ID     int64              `json:"id"`
	Option usecase.VoteOption `json:"option"`
}

type VotedEvent struct {
	ID     int64              `json:"id"`
	Option usecase.VoteOption `json:"option"`
	User   struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}
}

type FinishEvent struct {
	Result *domain.Place `json:"result"`
}

type ErrorEvent struct {
	Error string `json:"error"`
}

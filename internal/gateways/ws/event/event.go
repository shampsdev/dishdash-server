package event

import "dishdash.ru/internal/domain"

const (
	JoinLobby      = "joinLobby"
	UserJoined     = "userJoined"
	UserLeft       = "userLeft"
	SettingsUpdate = "settingsUpdate"
	StartSwipes    = "startSwipes"
	Place          = "card"
	Swipe          = "swipe"
	Match          = "match"
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
	PriceMin    int     `json:"priceMin"`
	PriceMax    int     `json:"priceMax"`
	MaxDistance int     `json:"maxDistance"`
	Tags        []int64 `json:"tags"`
}

type PlaceEvent struct {
	ID   int64         `json:"id"`
	Card *domain.Place `json:"card" mapstructure:"card"`
}

type SwipeEvent struct {
	SwipeType domain.SwipeType `json:"swipeType"`
}

type MatchEvent struct {
	ID   int           `json:"id"`
	Card *domain.Place `json:"card" mapstructure:"card"`
}

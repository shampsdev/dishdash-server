package swipes

import "dishdash.ru/internal/domain"

const (
	eventJoinLobby      = "joinLobby"
	eventSettingsUpdate = "settingsUpdate"
	eventCard           = "card"
	eventSwipe          = "swipe"
	eventMatch          = "match"
	eventUserJoined     = "userJoined"
	eventStartSwipes    = "startSwipes"
)

type swipeEvent struct {
	SwipeType domain.SwipeType `json:"swipeType"`
}

type matchEvent struct {
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
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

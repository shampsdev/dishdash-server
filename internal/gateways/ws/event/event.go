package event

import (
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
)

const (
	JoinLobby      = "joinLobby"
	LeaveLobby     = "leaveLobby"
	UserJoined     = "userJoined"
	UserLeft       = "userLeft"
	SettingsUpdate = "settingsUpdate"
	StartSwipes    = "startSwipes"
	Place          = "card"
	Swipe          = "swipe"
	Finish         = "finish"

	VoteAnnounce = "voteAnnounce"
	Vote         = "vote"
	Voted        = "voted"
	VoteResult   = "voteResult"

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

type VoteAnnounceEvent struct {
	usecase.Vote
}

type VoteEvent struct {
	VoteID   int64            `json:"voteId"`
	OptionID usecase.OptionID `json:"optionId"`
}

type VotedEvent struct {
	VoteID   int64            `json:"voteId"`
	OptionID usecase.OptionID `json:"optionId"`
	User     struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}
}

type VoteResultEvent struct {
	VoteID   int64            `json:"voteId"`
	OptionID usecase.OptionID `json:"optionId"`
}

type FinishEvent struct {
	Result  *domain.Place    `json:"result"`
	Matches []*usecase.Match `json:"matches"`
}

type ErrorEvent struct {
	Error string `json:"error"`
}

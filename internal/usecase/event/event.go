package event

import "dishdash.ru/internal/domain"

const (
	ErrorEvent = "error"

	JoinLobbyEvent      = "joinLobby"
	UserJoinedEvent     = "userJoined"
	LeaveLobbyEvent     = "leaveLobby"
	UserLeftEvent       = "userLeft"
	SettingsUpdateEvent = "settingsUpdate"
	StartSwipesEvent    = "startSwipes"
	CardsEvent          = "cards"
	SwipeEvent          = "swipe"
	MatchEvent          = "match"

	ResultsEvent = "results"
)

type Error struct {
	Error string `json:"error"`
}

func (e Error) Event() string { return ErrorEvent }

type JoinLobby struct {
	LobbyID string `json:"lobbyId"`
	UserID  string `json:"userId"`
}

func (e JoinLobby) Event() string { return JoinLobbyEvent }

type LeaveLobby struct{}

func (e LeaveLobby) Event() string { return LeaveLobbyEvent }

type UserJoined struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (e UserJoined) Event() string { return UserJoinedEvent }

type UserLeft struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (e UserLeft) Event() string { return UserLeftEvent }

type SettingsUpdate domain.LobbySettings

func (e SettingsUpdate) Event() string { return SettingsUpdateEvent }

type StartSwipes struct{}

func (e StartSwipes) Event() string { return StartSwipesEvent }

type Cards struct {
	Cards []*domain.Place `json:"cards"`
}

func (e Cards) Event() string { return CardsEvent }

type Swipe struct {
	SwipeType domain.SwipeType `json:"swipeType"`
}

func (e Swipe) Event() string { return SwipeEvent }

type Match struct {
	Card *domain.Place `json:"card"`
}

func (e Match) Event() string { return MatchEvent }

type Results struct {
	Top []TopPosition `json:"top"`
}

type TopPosition struct {
	Card  *domain.Place  `json:"card"`
	Likes []*domain.User `json:"likes"`
}

func (e Results) Event() string { return ResultsEvent }

package nevent

import "dishdash.ru/internal/domain"

const (
	ErrorEvent = "error"

	JoinLobbyEvent      = "joinLobby"
	LeaveLobbyEvent     = "leaveLobby"
	UserJoinedEvent     = "userJoined"
	UserLeftEvent       = "userLeft"
	SettingsUpdateEvent = "settingsUpdate"
	StartSwipesEvent    = "startSwipes"
	PlaceEvent          = "card"
	SwipeEvent          = "swipe"
	FinishEvent         = "finish"

	VoteAnnounceEvent = "voteAnnounce"
	VoteEvent         = "vote"
	VotedEvent        = "voted"
	VoteResultEvent   = "voteResult"
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

type SettingsUpdate struct {
	Location           domain.Coordinate          `json:"location"`
	UserID             string                     `json:"userId"`
	PriceMin           int                        `json:"priceMin"`
	PriceMax           int                        `json:"priceMax"`
	MaxDistance        int                        `json:"maxDistance"`
	Tags               []int64                    `json:"tags"`
	RecommendationOpts *domain.RecommendationOpts `json:"recommendation"`
}

func (e SettingsUpdate) Event() string { return SettingsUpdateEvent }

type StartSwipes struct{}

func (e StartSwipes) Event() string { return StartSwipesEvent }

type Place struct {
	ID   int64         `json:"id"`
	Card *domain.Place `json:"card"`
}

func (e Place) Event() string { return PlaceEvent }

type Swipe struct {
	SwipeType domain.SwipeType `json:"swipeType"`
}

func (e Swipe) Event() string { return SwipeEvent }

type VoteType string

const (
	VoteTypeMatch  VoteType = "match"
	VoteTypeFinish VoteType = "finish"
)

type OptionID int64

type VoteOption struct {
	ID   OptionID `json:"id"`
	Desc string   `json:"description"`
}

const (
	OptionIDLike OptionID = iota
	OptionIDDislike

	OptionIDFinish
	OptionIDContinue
)

type Match struct {
	ID    int           `json:"id"`
	Place *domain.Place `json:"card"`
}

type BaseVote struct {
	ID      int64               `json:"id"`
	Options []VoteOption        `json:"options"`
	Type    VoteType            `json:"type"`
	Votes   map[string]OptionID `json:"-"`
}

type VoteAnnounce interface {
	Event() string
	isVoteAnnounce()
}

type FinishVote struct {
	BaseVote
}

func (v FinishVote) Event() string   { return VoteAnnounceEvent }
func (v FinishVote) isVoteAnnounce() {}

type MatchVote struct {
	BaseVote
	Place *domain.Place `json:"card"`
}

func (v MatchVote) Event() string   { return VoteAnnounceEvent }
func (v MatchVote) isVoteAnnounce() {}

type VoteResult struct {
	Type     VoteType `json:"-"`
	VoteID   int64    `json:"voteId"`
	OptionID OptionID `json:"optionId"`
}

func (v VoteResult) Event() string { return VoteResultEvent }

type Vote struct {
	VoteID   int64    `json:"voteId"`
	OptionID OptionID `json:"optionId"`
}

func (v Vote) Event() string { return VoteEvent }

type Voted struct {
	VoteID   int64      `json:"voteId"`
	OptionID OptionID   `json:"optionId"`
	User     UserJoined `json:"user"`
}

func (v Voted) Event() string { return VotedEvent }

type Finish struct {
	Result  *domain.Place `json:"result"`
	Matches []*Match      `json:"matches"`
}

func (v Finish) Event() string { return FinishEvent }

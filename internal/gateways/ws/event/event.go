package event

const (
	JoinLobby  = "joinLobby"
	UserJoined = "userJoined"
	UserLeft   = "userLeft"
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

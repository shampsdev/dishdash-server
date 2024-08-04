package event

const (
	JoinLobby  = "joinLobby"
	UserJoined = "usersJoined"
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

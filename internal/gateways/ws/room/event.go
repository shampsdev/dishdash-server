package room

const (
	eventJoinLobby = "joinLobby"
)

type joinLobbyEvent struct {
	LobbyID string `json:"lobbyId"`
	UserID  string `json:"userId"`
}

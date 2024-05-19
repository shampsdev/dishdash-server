package dto

type Lobby struct {
	ID       int64  `json:"id"`
	Location string `json:"location"`
}

type LobbyToCreate struct {
	Location string `json:"location"`
}

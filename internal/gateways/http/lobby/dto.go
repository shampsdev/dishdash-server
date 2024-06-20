package lobby

import (
	"time"

	"dishdash.ru/internal/domain"
)

type lobbyOutput struct {
	ID        string            `json:"id"`
	CreatedAt time.Time         `json:"createdAt"`
	Location  domain.Coordinate `json:"location"`
}

type nearestLobbyOutput struct {
	Dist  float64     `json:"distance"`
	Lobby lobbyOutput `json:"lobby"`
}

type findLobbyInput struct {
	Dist     float64           `json:"dist"`
	Location domain.Coordinate `json:"location"`
}

func lobbyToOutput(lobby *domain.Lobby) lobbyOutput {
	return lobbyOutput{
		ID:        lobby.ID,
		CreatedAt: lobby.CreatedAt,
		Location:  lobby.Location,
	}
}

package lobby

import (
	"time"

	"dishdash.ru/internal/domain"
)

func ToLobbyOutput(lobby *domain.Lobby) *lobbyOutput {
	return &lobbyOutput{
		ID:        lobby.ID,
		State:     lobby.State,
		PriceAvg:  lobby.PriceAvg,
		Location:  lobby.Location,
		CreatedAt: lobby.CreatedAt,
		Tags:      lobby.Tags,
		Users:     lobby.Users,
	}
}

type lobbyOutput struct {
	ID        string            `json:"id"`
	State     domain.LobbyState `json:"state"`
	PriceAvg  int               `json:"priceAvg"`
	Location  domain.Coordinate `json:"location"`
	CreatedAt time.Time         `json:"createdAt"`

	Tags  []*domain.Tag  `json:"tags"`
	Users []*domain.User `json:"users"`
}

type nearestLobbyOutput struct {
	Dist  float64      `json:"distance"`
	Lobby *lobbyOutput `json:"lobby"`
}

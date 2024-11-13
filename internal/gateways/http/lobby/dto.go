package lobby

import (
	"dishdash.ru/internal/usecase"
)

type nearestLobbyOutput struct {
	Dist  float64      `json:"distance"`
	Lobby *usecase.LobbyOutput `json:"lobby"`
}

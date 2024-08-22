package lobby

import (
	"dishdash.ru/internal/domain"
)

type nearestLobbyOutput struct {
	Dist  float64       `json:"distance"`
	Lobby *domain.Lobby `json:"lobby"`
}

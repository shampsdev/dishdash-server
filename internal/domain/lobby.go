package domain

import (
	"time"

	"dishdash.ru/internal/dto"

	"github.com/kellydunn/golang-geo"
)

type Lobby struct {
	ID        int64
	CreatedAt time.Time
	Location  *geo.Point
}

func (lb *Lobby) ToDto() dto.Lobby {
	lobbyDto := dto.Lobby{
		ID:        lb.ID,
		CreatedAt: lb.CreatedAt,
	}

	lobbyDto.Location = Point2String(lb.Location)
	return lobbyDto
}

func (lb *Lobby) ParseDto(lobbyDto dto.Lobby) error {
	lb.ID = lobbyDto.ID
	lb.CreatedAt = lobbyDto.CreatedAt
	return ParsePoint(lobbyDto.Location, lb.Location)
}

func (lb *Lobby) ParseDtoToCreate(lobbyDto dto.LobbyToCreate) error {
	lb.CreatedAt = lobbyDto.CreatedAt
	return ParsePoint(lobbyDto.Location, lb.Location)
}

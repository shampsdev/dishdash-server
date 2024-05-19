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

func LobbyToDto(lobby Lobby) dto.Lobby {
	lobbyDto := dto.Lobby{
		ID: lobby.ID,
	}

	lobbyDto.Location = Point2String(lobby.Location)
	return lobbyDto
}

func LobbyFromDtoToCreate(lobbyDto dto.LobbyToCreate) (*Lobby, error) {
	lobby := &Lobby{
		Location: &geo.Point{},
	}

	var err error
	lobby.Location, err = ParsePoint(lobbyDto.Location)
	return lobby, err
}

func LobbyFromDto(lobbyDto dto.Lobby) (*Lobby, error) {
	lobby := &Lobby{
		ID:       lobbyDto.ID,
		Location: &geo.Point{},
	}

	var err error
	lobby.Location, err = ParsePoint(lobbyDto.Location)
	return lobby, err
}

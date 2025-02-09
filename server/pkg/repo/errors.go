package repo

import "errors"

var (
	ErrLobbyNotFound = errors.New("lobby not found")
	ErrPlaceNotFound = errors.New("place not found")
	ErrPlaceExists   = errors.New("place already exists")
)

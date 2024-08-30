package repo

import "errors"

var ErrLobbyNotFound = errors.New("lobby not found")
var ErrPlaceExists = errors.New("place already exists")

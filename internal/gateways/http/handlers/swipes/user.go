package swipes

import (
	"dishdash.ru/internal/domain"
	socketio "github.com/googollee/go-socket.io"
)

type user struct {
	ID     string
	lobby  *lobby
	swipes []swipe

	conn socketio.Conn
}

func (u *user) takeCard() *domain.Card {
	return u.lobby.takeCard(len(u.swipes))
}

func (u *user) swipe(s swipe) {
	u.swipes = append(u.swipes, s)
}

package newroom

import (
	"dishdash.ru/internal/gateways/ws/event"
	"dishdash.ru/internal/usecase"
	socketio "github.com/googollee/go-socket.io"
)

func SetupHandlers(sio *socketio.Server, cases usecase.Cases) {
	s := NewSocketIO(sio, cases)

	s.On(event.StartSwipes, (*usecase.Room).OnStartSwipes)
}

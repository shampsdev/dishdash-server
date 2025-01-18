package room

import (
	"dishdash.ru/internal/usecase"
	"dishdash.ru/internal/usecase/event"
	socketio "github.com/googollee/go-socket.io"
)

func SetupHandlers(sio *socketio.Server, cases usecase.Cases) {
	s := NewSocketIO(sio, cases)

	s.On(event.StartSwipesEvent,
		(*usecase.Room).OnStartSwipes)
	s.On(event.SwipeEvent,
		(*usecase.Room).OnSwipe)
	s.On(event.SettingsUpdateEvent,
		(*usecase.Room).OnSettingsUpdate)
	s.On(event.VoteEvent,
		(*usecase.Room).OnVote)
	s.On(event.LeaveLobbyEvent,
		(*usecase.Room).OnLeaveLobby)
}

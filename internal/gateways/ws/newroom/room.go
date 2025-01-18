package newroom

import (
	"dishdash.ru/internal/usecase"
	"dishdash.ru/internal/usecase/nevent"
	socketio "github.com/googollee/go-socket.io"
)

func SetupHandlers(sio *socketio.Server, cases usecase.Cases) {
	s := NewSocketIO(sio, cases)

	s.On(nevent.StartSwipesEvent,
		(*usecase.NRoom).OnStartSwipes)
	s.On(nevent.SwipeEvent,
		(*usecase.NRoom).OnSwipe)
	s.On(nevent.SettingsUpdateEvent,
		(*usecase.NRoom).OnSettingsUpdate)
	s.On(nevent.VoteEvent,
		(*usecase.NRoom).OnVote)
	s.On(nevent.LeaveLobbyEvent,
		(*usecase.NRoom).OnLeaveLobby)
}

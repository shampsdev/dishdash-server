package room

import (
	"dishdash.ru/pkg/usecase"
	"dishdash.ru/pkg/usecase/event"
	"dishdash.ru/pkg/usecase/state"
	socketio "github.com/googollee/go-socket.io"
)

func SetupHandlers(sio *socketio.Server, cases usecase.Cases) {
	s := NewSocketIO(sio, cases)

	s.On(event.StartSwipesEvent, state.WrapHMethod(
		(*usecase.Room).OnStartSwipes))
	s.On(event.SwipeEvent, state.WrapHMethod(
		(*usecase.Room).OnSwipe))
	s.On(event.SettingsUpdateEvent, state.WrapHMethod(
		(*usecase.Room).OnSettingsUpdate))
	s.On(event.LeaveLobbyEvent, state.WrapHMethod(
		(*usecase.Room).OnLeaveLobby))
}

package newroom

import (
	"context"
	"fmt"

	"dishdash.ru/internal/gateways/ws/event"
	"dishdash.ru/internal/usecase"
	"dishdash.ru/internal/usecase/nevent"
	"dishdash.ru/internal/usecase/state"

	socketio "github.com/googollee/go-socket.io"
	log "github.com/sirupsen/logrus"
)

type SocketIO struct {
	sio   *socketio.Server
	cases usecase.Cases
}

func NewSocketIO(sio *socketio.Server, cases usecase.Cases) *SocketIO {
	s := &SocketIO{
		sio:   sio,
		cases: cases,
	}
	s.setup()
	return s
}

func (s *SocketIO) setup() {
	s.sio.OnConnect("/", func(conn socketio.Conn) error {
		c := state.NewContext(s, wrapSocketIOConn(conn))
		conn.SetContext(c)
		return nil
	})

	s.sio.OnEvent("/", event.JoinLobby, func(conn socketio.Conn, joinEvent nevent.JoinLobby) {
		c, _ := conn.Context().(*state.Context[*usecase.NRoom])
		c.Log = log.WithFields(log.Fields{
			"user":  joinEvent.UserID,
			"room":  joinEvent.LobbyID,
			"event": nevent.JoinLobbyEvent,
		})

		user, err := s.cases.User.GetUserByID(context.Background(), joinEvent.UserID)
		if err != nil {
			c.Error(fmt.Errorf("error while getting user: %w", err))
			return
		}
		c.User = user

		room, err := s.cases.RoomRepo.GetRoom(context.Background(), joinEvent.LobbyID)
		if err != nil {
			c.Error(fmt.Errorf("error while getting room: %w", err))
			return
		}
		c.State = room
		c.Ctx = context.Background()

		err = c.State.OnJoin(c)
		if err != nil {
			c.Error(fmt.Errorf("error while adding user to room: %w", err))
			return
		}

		conn.Join(room.ID())
	})

	s.sio.OnError("/", func(conn socketio.Conn, err error) {
		c, ok := conn.Context().(*state.Context[*usecase.NRoom])
		if ok {
			c.Error(err)
		} else {
			log.Error(err)
			conn.Emit("error", err.Error())
		}
	})

	s.sio.OnDisconnect("/", func(conn socketio.Conn, msg string) {
		c, ok := conn.Context().(*state.Context[*usecase.NRoom])
		if !ok {
			log.Infof("disconnected with msg \"%s\" ", msg)
			return
		}

		if c.State == nil {
			return
		}

		err := c.State.OnLeave(c)
		if err != nil {
			c.Error(fmt.Errorf("error while removing user from room: %w", err))
		}
		if c.State.Empty() {
			err := s.cases.RoomRepo.DeleteRoom(context.Background(), c.State.ID())
			if err != nil {
				c.Error(fmt.Errorf("error while deleting room: %w", err))
			}
		}
		s.sio.LeaveRoom("/", c.State.ID(), conn)
	})
}

func (s *SocketIO) ForEach(roomID string, f func(c *state.Context[*usecase.NRoom])) {
	s.sio.ForEach("/", roomID, func(conn socketio.Conn) {
		f(conn.Context().(*state.Context[*usecase.NRoom]))
	})
}

func (s *SocketIO) On(event string, f interface{}) {
	s.sio.OnEvent("/", event, func(conn socketio.Conn, args interface{}) {
		c, _ := conn.Context().(*state.Context[*usecase.NRoom])
		if c.User == nil {
			c.Error(fmt.Errorf("not authenticated"))
			return
		}

		if c.State == nil {
			c.Error(fmt.Errorf("not in room"))
			return
		}

		c.Log = c.Log.WithFields(log.Fields{
			"room":  c.State.ID(),
			"user":  c.User.ID,
			"event": event,
		})
		c.Log.Debug("Event received")

		c.Ctx = context.Background()
		err := c.Call(f, args)
		if err != nil {
			c.Error(err)
		}
	})
}

type socketIOConn struct {
	conn socketio.Conn
}

func wrapSocketIOConn(conn socketio.Conn) *socketIOConn {
	return &socketIOConn{conn: conn}
}

func (s *socketIOConn) Emit(event string, data interface{}) {
	s.conn.Emit(event, data)
}

func (s *socketIOConn) Close() error {
	return s.conn.Close()
}

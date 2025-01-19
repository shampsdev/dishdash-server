package room

import (
	"context"
	"fmt"

	"dishdash.ru/internal/usecase"
	"dishdash.ru/internal/usecase/event"
	"dishdash.ru/internal/usecase/state"

	socketio "github.com/googollee/go-socket.io"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type SocketIO struct {
	sio     *socketio.Server
	metrics ServerMetrics
	cases   usecase.Cases
}

func NewSocketIO(sio *socketio.Server, cases usecase.Cases) *SocketIO {
	s := &SocketIO{
		sio:     sio,
		cases:   cases,
		metrics: NewServerMetrics(),
	}
	s.setup()
	return s
}

func (s *SocketIO) setup() {
	s.sio.OnConnect("/", func(conn socketio.Conn) error {
		defer s.metrics.ActiveConnections.Inc()
		c := state.NewContext(s, s.wrapSocketIOConn(conn))
		conn.SetContext(c)
		return nil
	})

	s.sio.OnEvent("/", event.JoinLobbyEvent, func(conn socketio.Conn, joinEvent event.JoinLobby) {
		c, _ := conn.Context().(*state.Context[*usecase.Room])
		c.Log = log.WithFields(log.Fields{
			"user":  joinEvent.UserID,
			"room":  joinEvent.LobbyID,
			"event": event.JoinLobbyEvent,
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

		c.Log.Debug("Event received")

		err = c.State.OnJoin(c)
		if err != nil {
			c.Error(fmt.Errorf("error while adding user to room: %w", err))
			return
		}

		conn.Join(room.ID())
	})

	s.sio.OnError("/", func(conn socketio.Conn, err error) {
		c, ok := conn.Context().(*state.Context[*usecase.Room])
		if ok {
			c.Error(err)
		} else {
			log.Error(err)
			conn.Emit("error", err.Error())
		}
	})

	s.sio.OnDisconnect("/", func(conn socketio.Conn, msg string) {
		defer func() {
			s.metrics.ActiveConnections.Dec()
		}()
		c, ok := conn.Context().(*state.Context[*usecase.Room])
		if !ok {
			log.Infof("disconnected with msg \"%s\" ", msg)
			return
		}

		if c.State == nil {
			return
		}

		c.Log = log.WithFields(log.Fields{
			"user": c.User.ID,
			"room": c.State.ID(),
		})

		err := c.State.OnLeave(c)
		if err != nil {
			c.Error(fmt.Errorf("error while removing user from room: %w", err))
		}
		c.Log.Info("Leave room")
		s.sio.LeaveRoom("/", c.State.ID(), conn)
	})
}

func (s *SocketIO) ForEach(roomID string, f func(c *state.Context[*usecase.Room])) {
	s.sio.ForEach("/", roomID, func(conn socketio.Conn) {
		f(conn.Context().(*state.Context[*usecase.Room]))
	})
}

func (s *SocketIO) On(event string, f state.HandlerFunc[*usecase.Room]) {
	s.sio.OnEvent("/", event, func(conn socketio.Conn, arg interface{}) {
		s.metrics.Requests.WithLabelValues(event).Inc()
		timer := prometheus.NewTimer(s.metrics.ResponseDuration.WithLabelValues(event))
		defer timer.ObserveDuration()

		c, _ := conn.Context().(*state.Context[*usecase.Room])
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
		err := f(c, arg)
		if err != nil {
			c.Error(err)
		}
	})
}

type socketIOConn struct {
	conn socketio.Conn
	sio  *SocketIO
}

func (s *SocketIO) wrapSocketIOConn(conn socketio.Conn) *socketIOConn {
	return &socketIOConn{conn: conn, sio: s}
}

func (s *socketIOConn) Emit(event string, data interface{}) {
	s.sio.metrics.Responses.WithLabelValues(event).Inc()
	s.conn.Emit(event, data)
}

func (s *socketIOConn) Close() error {
	return s.conn.Close()
}

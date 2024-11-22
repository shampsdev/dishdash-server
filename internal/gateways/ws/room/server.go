package room

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sync"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/gateways/ws/event"
	"dishdash.ru/internal/usecase"
	socketio "github.com/googollee/go-socket.io"
	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	SIO *socketio.Server

	Metrics ServerMetrics
}

type Context struct {
	lock sync.RWMutex

	Conn   socketio.Conn
	Server *Server
	Log    *log.Entry

	User *domain.User
	Room *usecase.Room
}

type EventOpts struct {
	// Allowed lobby states
	// If empty, then any lobby state is allowed
	Allowed []domain.LobbyState
}

func NewServer(sio *socketio.Server) *Server {
	s := &Server{
		SIO:     sio,
		Metrics: NewServerMetrics(),
	}

	sio.OnConnect("/", func(conn socketio.Conn) error {
		defer s.Metrics.ActiveConnections.Inc()

		log.Println("connected: ", conn.ID())
		conn.SetContext(&Context{
			Conn:   conn,
			Server: s,
			Log:    log.WithField("user", conn.ID()),
		})
		return nil
	})

	return s
}

func (c *Context) HandleError(err error) {
	c.Log.Error(err)
	c.Conn.Emit(event.Error, event.ErrorEvent{
		Error: err.Error(),
	})
	if err := c.Conn.Close(); err != nil {
		c.Log.WithError(err).Error("Error closing connection")
	}
}

func (c *Context) Emit(eventName string, args ...interface{}) {
	c.Log.Debugf("Emit %s", eventName)
	c.Conn.Emit(eventName, args...)
	c.Server.Metrics.Responses.WithLabelValues(eventName).Inc()
}

func (s *Server) GetContext(conn socketio.Conn) (*Context, bool) {
	if conn.Context() == nil {
		return nil, false
	}
	c, ok := conn.Context().(*Context)
	if !ok {
		err := errors.New("context not found (maybe did not join)")
		log.Error(err)
		conn.Emit(event.Error, event.ErrorEvent{
			Error: err.Error(),
		})
		if err := conn.Close(); err != nil {
			log.WithError(err).Error("Error closing connection")
		}
		return nil, false
	}

	if c.Room != nil && c.User != nil {
		c.Log = log.WithFields(log.Fields{
			"room": c.Room.ID,
			"user": c.User.ID,
		})
	}
	return c, true
}

func (s *Server) On(
	eventName string,
	opts EventOpts,
	f interface{},
) {
	s.SIO.OnEvent("/", eventName, func(conn socketio.Conn, eventData interface{}) {
		s.Metrics.Requests.WithLabelValues(eventName).Inc()
		timer := prometheus.NewTimer(s.Metrics.ResponseDuration.WithLabelValues(eventName))
		defer timer.ObserveDuration()

		c, ok := s.GetContext(conn)
		if !ok {
			return
		}
		c.Log = c.Log.WithField("event", eventName)
		c.Log.Debug("Event received")

		if c.Room != nil {
			if len(opts.Allowed) != 0 && !slices.Contains(opts.Allowed, c.Room.State()) {
				c.HandleError(fmt.Errorf("event '%s' is not allowed in lobby state '%s'", eventName, c.Room.State()))
				return
			}
		}

		// Black magic
		if reflect.TypeOf(f).NumIn() == 1 {
			reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(c)})
			return
		}

		eventDataCasted := reflect.New(reflect.TypeOf(f).In(1)).Elem().Interface()
		err := mapstructure.Decode(eventData, &eventDataCasted)
		if err != nil {
			c.HandleError(err)
			return
		}
		c.Conn = conn
		reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(c), reflect.ValueOf(eventDataCasted)})
	})
}

package state

import (
	"context"

	"dishdash.ru/pkg/domain"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type Context[State WithID] struct {
	Ctx context.Context

	State State
	User  *domain.User
	Log   *log.Entry

	s    Server[State]
	conn Conn
}

func NewContext[State WithID](s Server[State], conn Conn) *Context[State] {
	return &Context[State]{
		s:    s,
		conn: conn,
		Log:  logrus.NewEntry(logrus.New()),
	}
}

type Event interface {
	Event() string
}

func (c *Context[State]) Emit(e Event) {
	c.conn.Emit(e.Event(), e)

	log.WithFields(log.Fields{
		"event": e.Event(),
		"user":  c.User.ID,
		"room":  c.State.ID(),
	}).Debug("Event emitted")
}

func (c *Context[State]) Broadcast(e Event) {
	c.s.ForEach(c.State.ID(), func(c *Context[State]) {
		c.Emit(e)
	})
}

func (c *Context[State]) BroadcastToOthers(e Event) {
	c.s.ForEach(c.State.ID(), func(cc *Context[State]) {
		if cc.User.ID != c.User.ID {
			cc.Emit(e)
		}
	})
}

func (c *Context[State]) ForEach(f func(c *Context[State])) {
	c.s.ForEach(c.State.ID(), f)
}

func (c *Context[State]) Close() error {
	return c.conn.Close()
}

func (c *Context[State]) Error(err error) {
	c.Log.Error(err)
	c.conn.Emit("error", err.Error())
	if err = c.conn.Close(); err != nil {
		c.Log.Error(err)
	}
}

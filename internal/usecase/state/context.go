package state

import (
	"context"

	"dishdash.ru/internal/domain"
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

func NewContext[State WithID](s Server[State], conn Conn) Context[State] {
	return Context[State]{
		s:    s,
		conn: conn,
		Log:  logrus.NewEntry(logrus.New()),
	}
}

func (c *Context[State]) Emit(event string, data interface{}) {
	c.conn.Emit(event, data)
}

func (c *Context[State]) Broadcast(event string, data interface{}) {
	c.s.ForEach(c.State.ID(), func(c *Context[State]) {
		c.Emit(event, data)
	})
}

func (c *Context[State]) BroadcastToOthers(event string, data interface{}) {
	c.s.ForEach(c.State.ID(), func(cc *Context[State]) {
		if cc.User.ID != c.User.ID {
			cc.Emit(event, data)
		}
	})
}

func (c *Context[State]) ForEach(f func(c *Context[State])) {
	c.s.ForEach(c.State.ID(), f)
}

func (c *Context[State]) Close() {
	c.conn.Close()
}

func (c *Context[State]) Error(err error) {
	c.Log.Error(err)
	c.conn.Emit("error", err.Error())
	if err = c.conn.Close(); err != nil {
		c.Log.Error(err)
	}
}

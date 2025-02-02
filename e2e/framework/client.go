package framework

import (
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase/event"
	socketio "github.com/googollee/go-socket.io"
	log "github.com/sirupsen/logrus"
)

type Event interface {
	Event() string
}

type HandlerFunc func(c *Client, arg interface{}) error

type Client struct {
	User *domain.User
	Log  *log.Entry

	fw  *Framework
	cli *socketio.Client
}

func (c *Client) On(event string, f HandlerFunc) {
	c.cli.OnEvent(event, func(_ socketio.Conn, arg interface{}) {
		c.Log = c.Log.WithFields(log.Fields{
			"user":  c.User.ID,
			"event": event,
		})
		c.Log.Debug("Event received")
		c.fw.Session.RecordEvent(c.User, event, arg)
		err := f(c, arg)
		if err != nil {
			c.Log.Error(err)
		}
	})
}

func (c *Client) setup(toRecord map[string]struct{}) {
	c.cli.OnConnect(func(_ socketio.Conn) error {
		c.Log.Debug("Connected")
		return nil
	})
	c.cli.OnDisconnect(func(_ socketio.Conn, reason string) {
		c.Log.Debugf("Disconnected: %s", reason)
	})
	c.cli.OnError(func(_ socketio.Conn, err error) {
		c.Log.Error(err)
	})

	for ev := range toRecord {
		c.On(ev, func(_ *Client, _ interface{}) error { return nil })
	}
}

func (c *Client) JoinLobby(lobby *domain.Lobby) {
	c.Emit(event.JoinLobby{
		LobbyID: lobby.ID,
		UserID:  c.User.ID,
	})
}

func (c *Client) Emit(ev Event) {
	c.cli.Emit(ev.Event(), ev)

	c.Log.WithFields(log.Fields{
		"user":  c.User.ID,
		"event": ev.Event(),
	}).Debug("Event emitted")
}

func (c *Client) Close() error {
	return c.cli.Close()
}

package room

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/gateways/ws/event"
	"dishdash.ru/internal/usecase"

	socketio "github.com/googollee/go-socket.io"
)

type Context struct {
	lock sync.RWMutex
	User *domain.User
	Room *usecase.Room
	Log  *log.Entry
}

func SetupHandlers(s *socketio.Server, cases usecase.Cases) {
	handleError := func(conn socketio.Conn, err error) {
		if c, ok := conn.Context().(*Context); ok {
			c.Log.Error(err)
		} else {
			log.Error(err)
		}

		conn.Emit(event.Error, event.ErrorEvent{
			Error: err.Error(),
		})
		if err := conn.Close(); err != nil {
			log.WithError(err).Error("Error closing connection")
		}
	}

	getContext := func(conn socketio.Conn) (*Context, bool) {
		c, ok := conn.Context().(*Context)
		if !ok {
			handleError(conn, fmt.Errorf("context not found"))
			return nil, false
		}
		return c, ok
	}

	s.OnConnect("/", func(s socketio.Conn) error {
		log.Println("connected: ", s.ID())
		s.SetContext("")
		return nil
	})

	s.OnEvent("/", event.JoinLobby, func(conn socketio.Conn, joinEvent event.JoinLobbyEvent) {
		user, err := cases.User.GetUserByID(context.Background(), joinEvent.UserID)
		if err != nil {
			handleError(conn, fmt.Errorf("error while getting user: %w", err))
			return
		}

		room, err := cases.RoomRepo.GetRoom(context.Background(), joinEvent.LobbyID)
		if err != nil {
			handleError(conn, fmt.Errorf("error while getting room: %w", err))
			return
		}

		err = room.AddUser(user)
		if err != nil {
			handleError(conn, fmt.Errorf("error while adding user to room: %w", err))
			return
		}

		conn.Join(room.ID)
		c := &Context{
			User: user,
			Room: room,
			Log: log.WithFields(log.Fields{
				"room": room.ID,
				"user": user.ID,
			}),
		}
		c.lock.Lock()
		conn.SetContext(c)
		c.lock.Unlock()
		c.Log.Info("User joined")

		broadcastToOthersInRoom(s, c.User.ID, c.Room.ID, event.UserJoined,
			event.UserJoinedEvent{
				ID:     user.ID,
				Name:   user.Name,
				Avatar: user.Avatar,
			})

		for _, u := range c.Room.Users() {
			conn.Emit(event.UserJoined, event.UserJoinedEvent{
				ID:     u.ID,
				Name:   u.Name,
				Avatar: u.Avatar,
			})
		}
		settings := c.Room.Settings()
		conn.Emit(event.SettingsUpdate, event.SettingsUpdateEvent{
			Location:    settings.Location,
			PriceMin:    settings.PriceAvg - 300,
			PriceMax:    settings.PriceAvg + 300,
			MaxDistance: 4000,
			Tags:        settings.Tags,
		})

		if c.Room.Swiping() {
			conn.Emit(event.StartSwipes)
		}

		if room.Finished() {
			conn.Emit(event.Finish, event.FinishEvent{
				Result:  room.Result(),
				Matches: room.Matches(),
			})
			return
		}
	})

	s.OnEvent("/", event.SettingsUpdate, func(conn socketio.Conn, se event.SettingsUpdateEvent) {
		c, ok := getContext(conn)
		if !ok {
			return
		}

		ctx := context.Background()
		err := c.Room.UpdateLobbySettings(ctx, se.Location, (se.PriceMax+se.PriceMax)/2, se.Tags, nil)
		if err != nil {
			handleError(conn, fmt.Errorf("error while updating lobby: %w", err))
			return
		}

		se.UserID = c.User.ID
		s.BroadcastToRoom("", c.Room.ID, event.SettingsUpdate, se)
	})

	s.OnEvent("/", event.StartSwipes, func(conn socketio.Conn) {
		c, ok := getContext(conn)
		if !ok {
			return
		}

		err := c.Room.StartSwipes(context.Background())
		if err != nil {
			handleError(conn, fmt.Errorf("error while starting swipes: %w", err))
			return
		}

		c.Log.Info("Start swipes")
		s.BroadcastToRoom("", c.Room.ID, event.StartSwipes)
		s.ForEach("", c.Room.ID, func(conn socketio.Conn) {
			c, ok := getContext(conn)
			if !ok {
				return
			}

			p := c.Room.GetNextPlaceForUser(c.User.ID)
			conn.Emit(event.Place, event.PlaceEvent{
				ID:   p.ID,
				Card: p,
			})
		})
	})

	s.OnEvent("/", event.Swipe, func(conn socketio.Conn, se event.SwipeEvent) {
		c, ok := getContext(conn)
		if !ok {
			return
		}
		m, err := c.Room.Swipe(c.User.ID, c.Room.GetNextPlaceForUser(c.User.ID).ID, se.SwipeType)
		if err != nil {
			handleError(conn, fmt.Errorf("error while swiping: %w", err))
			return
		}

		if m != nil {
			s.BroadcastToRoom("/", c.Room.ID, event.Match,
				event.MatchEvent{
					ID:   m.ID,
					Card: m.Place,
				})
		}
		p := c.Room.GetNextPlaceForUser(c.User.ID)
		conn.Emit(event.Place, event.PlaceEvent{
			ID:   p.ID,
			Card: p,
		})
	})

	s.OnError("/", func(conn socketio.Conn, e error) {
		handleError(conn, e)
	})

	voteLock := sync.RWMutex{}
	s.OnEvent("/", event.Vote, func(conn socketio.Conn, ve event.VoteEvent) {
		voteLock.Lock()
		defer voteLock.Unlock()
		c, ok := getContext(conn)
		if !ok {
			return
		}
		err := c.Room.Vote(c.User.ID, ve.Option)
		if err != nil {
			handleError(conn, fmt.Errorf("error while voting: %w", err))
			return
		}
		s.BroadcastToRoom("/", c.Room.ID, event.Voted, event.VotedEvent{
			ID:     ve.ID,
			Option: ve.Option,
			User: struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Avatar string `json:"avatar"`
			}{
				ID:     c.User.ID,
				Name:   c.User.Name,
				Avatar: c.User.Avatar,
			},
		})

		if c.Room.Swiping() || c.Room.Finished() {
			s.BroadcastToRoom("/", c.Room.ID, event.ReleaseMatch)
		}
		if c.Room.Finished() {
			s.BroadcastToRoom("/", c.Room.ID, event.Finish, event.FinishEvent{
				Result:  c.Room.Result(),
				Matches: c.Room.Matches(),
			})
		}
	})

	s.OnEvent("/", event.LeaveLobby, func(conn socketio.Conn) {
		if err := conn.Close(); err != nil {
			log.WithError(err).Error("error while closing connection")
		}
	})

	s.OnDisconnect("/", func(conn socketio.Conn, msg string) {
		c, ok := getContext(conn)
		if !ok {
			log.Infof("disconnected with msg \"%s\" ", msg)
			return
		}

		err := c.Room.RemoveUser(c.User.ID)
		if err != nil {
			handleError(conn, fmt.Errorf("error while removing user from room: %w", err))
			return
		}

		s.LeaveRoom("", c.Room.ID, conn)
		s.BroadcastToRoom("", c.Room.ID, event.UserLeft,
			event.UserLeftEvent{
				ID:     c.User.ID,
				Name:   c.User.Name,
				Avatar: c.User.Avatar,
			})

		c.Log.Info("Leave room")
		if c.Room.Empty() {
			err = cases.RoomRepo.DeleteRoom(context.Background(), c.Room.ID)
			if err != nil {
				c.Log.WithError(err).Error("error while deleting room")
			}
		}
	})
}

func broadcastToOthersInRoom(s *socketio.Server, userID, room, event string, args ...interface{}) {
	s.ForEach("", room, func(conn socketio.Conn) {
		c, ok := conn.Context().(*Context)
		if !ok {
			return
		}
		c.lock.RLock()
		defer c.lock.RUnlock()
		if c.User.ID != userID {
			conn.Emit(event, args...)
		}
	})
}

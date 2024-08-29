package room

import (
	"context"
	"log"
	"sync"

	"dishdash.ru/internal/gateways/ws/event"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"

	socketio "github.com/googollee/go-socket.io"
)

type Context struct {
	lock sync.RWMutex
	User *domain.User
	Room *usecase.Room
}

func SetupHandlers(s *socketio.Server, cases usecase.Cases) {
	s.OnConnect("/", func(s socketio.Conn) error {
		log.Println("connected: ", s.ID())
		s.SetContext("")
		return nil
	})

	s.OnEvent("/", event.JoinLobby, func(conn socketio.Conn, joinEvent event.JoinLobbyEvent) {
		user, err := cases.User.GetUserByID(context.Background(), joinEvent.UserID)
		if err != nil {
			log.Println("error while getting user: ", err)
			_ = conn.Close()
			return
		}

		room, err := cases.RoomRepo.GetRoom(context.Background(), joinEvent.LobbyID)
		if err != nil {
			log.Println("error while getting room: ", err)
			_ = conn.Close()
			return
		}

		err = room.AddUser(user)
		if err != nil {
			log.Println("error while adding user to room: ", err)
			_ = conn.Close()
			return
		}

		conn.Join(room.ID)
		c := &Context{
			User: user,
			Room: room,
		}
		c.lock.Lock()
		conn.SetContext(c)
		c.lock.Unlock()
		log.Printf("<user %s> joined to <lobby %s>", joinEvent.UserID, joinEvent.LobbyID)

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
			PriceMin:    settings.PriceAvg - 300,
			PriceMax:    settings.PriceAvg + 300,
			MaxDistance: 4000,
			Tags:        settings.Tags,
		})

		if c.Room.Swiping() {
			conn.Emit(event.StartSwipes)
		}
	})

	s.OnEvent("/", event.SettingsUpdate, func(conn socketio.Conn, se event.SettingsUpdateEvent) {
		c, ok := conn.Context().(*Context)
		if !ok {
			_ = conn.Close()
			return
		}

		ctx := context.Background()
		err := c.Room.UpdateLobby(ctx, (se.PriceMax+se.PriceMax)/2, se.Tags, nil)
		if err != nil {
			log.Println("error while updating lobby: ", err)
			_ = conn.Close()
			return
		}

		s.BroadcastToRoom("", c.Room.ID, event.SettingsUpdate, se)
	})

	s.OnEvent("/", event.StartSwipes, func(conn socketio.Conn) {
		c, ok := conn.Context().(*Context)
		if !ok {
			_ = conn.Close()
			return
		}

		err := c.Room.StartSwipes(context.Background())
		if err != nil {
			log.Println("error while starting swipes: ", err)
			_ = conn.Close()
		}

		s.ForEach("", c.Room.ID, func(conn socketio.Conn) {
			c, ok := conn.Context().(*Context)
			if !ok {
				_ = conn.Close()
			}

			p := c.Room.GetNextPlaceForUser(c.User.ID)
			conn.Emit(event.Place, event.PlaceEvent{
				ID:   p.ID,
				Card: p,
			})
		})
		log.Printf("start swipes in <lobby %s>", c.Room.ID)
	})

	s.OnEvent("/", event.Swipe, func(conn socketio.Conn, se event.SwipeEvent) {
		c, ok := conn.Context().(*Context)
		if !ok {
			_ = conn.Close()
			return
		}
		m, err := c.Room.Swipe(c.User.ID, c.Room.GetNextPlaceForUser(c.User.ID).ID, se.SwipeType)
		if err != nil {
			log.Println("error while swiping: ", err)
			_ = conn.Close()
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

	s.OnError("/", func(_ socketio.Conn, e error) {
		log.Println("faced error: ", e)
	})

	voteLock := sync.RWMutex{}
	s.OnEvent("/", event.Vote, func(conn socketio.Conn, ve event.VoteEvent) {
		voteLock.Lock()
		defer voteLock.Unlock()
		c, ok := conn.Context().(*Context)
		if !ok {
			_ = conn.Close()
			return
		}
		c.Room.Vote(c.User.ID, ve.Option)
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
				Result: c.Room.Result(),
			})
		}
	})

	s.OnDisconnect("/", func(conn socketio.Conn, msg string) {
		c, ok := conn.Context().(*Context)
		if !ok {
			log.Println("disconnected: ", msg)
			_ = conn.Close()
			return
		}

		err := c.Room.RemoveUser(c.User.ID)
		if err != nil {
			log.Println("error while removing user from room: ", err)
			_ = conn.Close()
			return
		}

		s.LeaveRoom("", c.Room.ID, conn)
		s.BroadcastToRoom("", c.Room.ID, event.UserLeft,
			event.UserLeftEvent{
				ID:     c.User.ID,
				Name:   c.User.Name,
				Avatar: c.User.Avatar,
			})

		log.Printf("<user %s> leave <lobby %s>", c.User.ID, c.Room.ID)
		if c.Room.Empty() {
			err = cases.RoomRepo.DeleteRoom(context.Background(), c.Room.ID)
			if err != nil {
				log.Println("error while deleting room: ", err)
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

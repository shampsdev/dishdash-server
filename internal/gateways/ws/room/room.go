package room

import (
	"context"
	"log"

	"dishdash.ru/internal/gateways/ws/event"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"

	socketio "github.com/googollee/go-socket.io"
)

type Context struct {
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

		conn.Join(room.Lobby.ID)
		broadcastToOthersInRoom(
			s, user.ID, room.Lobby.ID, event.UserJoined,
			event.UserJoinedEvent{
				ID:     user.ID,
				Name:   user.Name,
				Avatar: user.Avatar,
			},
		)
		conn.SetContext(Context{
			User: user,
			Room: room,
		})

		log.Printf("<user %s> joined to <lobby %s>", joinEvent.UserID, joinEvent.LobbyID)
	})

	s.OnEvent("/", event.SettingsUpdate, func(conn socketio.Conn, se event.SettingsUpdateEvent) {
		c, ok := conn.Context().(Context)
		if !ok {
			_ = conn.Close()
			return
		}

		ctx := context.Background()
		err := c.Room.UpdateLobby(ctx, usecase.UpdateLobbyInput{
			ID: c.Room.Lobby.ID,
			SaveLobbyInput: usecase.SaveLobbyInput{
				PriceAvg: (se.PriceMin + se.PriceMax) / 2,
				Location: c.Room.Lobby.Location,
				Tags:     se.Tags,
			},
		})
		if err != nil {
			log.Println("error while updating lobby: ", err)
			_ = conn.Close()
			return
		}
	})

	s.OnEvent("/", event.StartSwipes, func(conn socketio.Conn) {
		c, ok := conn.Context().(Context)
		if !ok {
			_ = conn.Close()
			return
		}

		err := c.Room.StartSwipes(context.Background())
		if err != nil {
			log.Println("error while starting swipes: ", err)
			_ = conn.Close()
		}

		c.Room.UsersPlaceMutex.Lock()
		s.ForEach("", c.Room.Lobby.ID, func(conn socketio.Conn) {
			c, ok := conn.Context().(Context)
			if !ok {
				_ = conn.Close()
			}
			p := c.Room.UsersPlace[c.User.ID]
			conn.Emit(event.Place, event.PlaceEvent{
				ID:   p.ID,
				Card: p,
			})
		})
		c.Room.UsersPlaceMutex.Unlock()
		log.Printf("start swipes in <lobby %s>", c.Room.Lobby.ID)
	})

	s.OnEvent("/", event.Swipe, func(conn socketio.Conn, se event.SwipeEvent) {
		c, ok := conn.Context().(Context)
		if !ok {
			_ = conn.Close()
			return
		}
		m, err := c.Room.Swipe(c.User.ID, c.Room.UsersPlace[c.User.ID].ID, se.SwipeType)
		if err != nil {
			log.Println("error while swiping: ", err)
			_ = conn.Close()
			return
		}

		if m != nil {
			s.BroadcastToRoom("/", c.Room.Lobby.ID, event.Match,
				event.MatchEvent{
					ID:   m.ID,
					Card: m.Place,
				})
		}
		p := c.Room.UsersPlace[c.User.ID]
		conn.Emit(event.Place, event.PlaceEvent{
			ID:   p.ID,
			Card: p,
		})
	})

	s.OnError("/", func(_ socketio.Conn, e error) {
		log.Println("faced error: ", e)
	})

	s.OnDisconnect("/", func(conn socketio.Conn, msg string) {
		c, ok := conn.Context().(Context)
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

		broadcastToOthersInRoom(s, c.User.ID, c.Room.Lobby.ID, event.UserLeft,
			event.UserLeftEvent{
				ID:     c.User.ID,
				Name:   c.User.Name,
				Avatar: c.User.Avatar,
			},
		)

		log.Printf("<user %s> leave <lobby %s>", c.User.ID, c.Room.Lobby.ID)
		if len(c.Room.Users) == 0 {
			err = cases.RoomRepo.DeleteRoom(context.Background(), c.Room.Lobby.ID)
			if err != nil {
				log.Println("error while deleting room: ", err)
			}
		}
	})
}

func broadcastToOthersInRoom(s *socketio.Server, userID, room, event string, args ...interface{}) {
	s.ForEach("", room, func(conn socketio.Conn) {
		c, ok := conn.Context().(Context)
		if !ok {
			return
		}
		if c.User.ID != userID {
			conn.Emit(event, args...)
		}
	})
}
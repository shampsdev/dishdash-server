package room

import (
	"context"
	"encoding/json"
	"log"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
	"dishdash.ru/internal/usecase/room"

	socketio "github.com/googollee/go-socket.io"
)

type Context struct {
	User *domain.User
	Room *room.Room
}

func SetupHandlers(s *socketio.Server, useCases usecase.Cases, roomRepo room.Repo) {
	s.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})

	s.OnEvent("/", eventJoinLobby, func(conn socketio.Conn, msg string) {
		var joinEvent joinLobbyEvent
		err := json.Unmarshal([]byte(msg), &joinEvent)
		if err != nil {
			_ = conn.Close()
			return
		}

		user, err := useCases.User.GetUserByID(context.Background(), joinEvent.UserID)
		if err != nil {
			_ = conn.Close()
			return
		}

		room, err := roomRepo.GetRoom(context.Background(), joinEvent.LobbyID)
		if err != nil {
			_ = conn.Close()
			return
		}

		err = room.AddUser(user)
		if err != nil {
			_ = conn.Close()
			return
		}

		conn.Join(room.Lobby.ID)
		conn.SetContext(Context{
			User: user,
			Room: room,
		})

		log.Printf("<user %s> joined to <lobby %s>", joinEvent.UserID, joinEvent.LobbyID)
	})

	s.OnDisconnect("/", func(conn socketio.Conn, _ string) {
		c, ok := conn.Context().(Context)
		if !ok {
			log.Println("user not registered, disconnected on disconnect")
			_ = conn.Close()
			return
		}

		err := c.Room.RemoveUser(c.User.ID)
		if err != nil {
			_ = conn.Close()
			return
		}

		log.Printf("<user %s> leave <lobby %s>", c.User.ID, c.Room.Lobby.ID)
		if len(c.Room.Users) == 0 {
			_ = roomRepo.DeleteRoom(context.Background(), c.Room.Lobby.ID)
		}
	})
}

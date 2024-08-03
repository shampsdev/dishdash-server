package room

import (
	"context"
	"encoding/json"
	"log"

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

	s.OnEvent("/", eventJoinLobby, func(conn socketio.Conn, msg string) {
		var joinEvent joinLobbyEvent
		err := json.Unmarshal([]byte(msg), &joinEvent)
		if err != nil {
			log.Println("error while unmarshalling join lobby: ", err)
			_ = conn.Close()
			return
		}

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
		conn.SetContext(Context{
			User: user,
			Room: room,
		})

		log.Printf("<user %s> joined to <lobby %s>", joinEvent.UserID, joinEvent.LobbyID)
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

		log.Printf("<user %s> leave <lobby %s>", c.User.ID, c.Room.Lobby.ID)
		if len(c.Room.Users) == 0 {
			err = cases.RoomRepo.DeleteRoom(context.Background(), c.Room.Lobby.ID)
			if err != nil {
				log.Println("error while deleting room: ", err)
			}
		}
	})
}

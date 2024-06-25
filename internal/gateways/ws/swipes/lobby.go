package swipes

import (
	"context"
	"encoding/json"
	"log"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/entities"
	"dishdash.ru/internal/usecase"

	socketio "github.com/googollee/go-socket.io"
)

func SetupHandlers(s *socketio.Server, useCases usecase.Cases) {
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
		_, ok := conn.Context().(*domain.User)
		if ok {
			return
		}

		domLobby, err := useCases.Lobby.GetLobbyByID(context.Background(), joinEvent.LobbyID)
		if err != nil {
			_ = conn.Close()
			return
		}

		lobby, err := entities.FindLobby(domLobby, useCases.Card)
		if err != nil {
			_ = conn.Close()
			return
		}

		u, err := useCases.User.GetUserByID(context.Background(), joinEvent.UserID)
		if err != nil {
			_ = conn.Close()
			return
		}

		user := entities.NewUser(*u, useCases.Swipe)
		lobby.Register(conn.ID(), user)

		conn.Join(domLobby.ID)
		conn.SetContext(user)

		s.BroadcastToRoom(
			"",
			user.Lobby.Id,
			"userJoined",
			userJoinEvent{
				Name:   u.Name,
				Avatar: u.Avatar,
			},
		)

		firstCard := user.Card()
		conn.Emit(eventCard, cardEvent{Card: *firstCard})
	})

	s.OnEvent("/", eventSettingsUpdate, func(conn socketio.Conn, msg string) {
		var updateEvent settingsUpdateEvent
		err := json.Unmarshal([]byte(msg), &updateEvent)
		if err != nil {
			log.Println("wrong settings change event")
			_ = conn.Close()
			return
		}

		user, ok := conn.Context().(*entities.User)
		if !ok {
			log.Println("user not registered, disconnected")
			_ = conn.Close()
			return
		}

		s.BroadcastToRoom(
			"",
			user.Lobby.Id,
			"settingsUpdate",
			updateEvent,
		)
	})

	s.OnEvent("/", eventSwipe, func(conn socketio.Conn, msg string) {
		var swipeEvent swipeEvent
		err := json.Unmarshal([]byte(msg), &swipeEvent)
		if err != nil {
			log.Println("wrong swipe event")
			_ = conn.Close()
			return
		}

		u, ok := conn.Context().(*entities.User)
		if !ok {
			log.Println("user not registered, disconnected")
			_ = conn.Close()
			return
		}

		card := u.Card()
		log.Println(card)

		match := u.Swipe(swipeEvent.SwipeType)
		if match != nil {
			s.BroadcastToRoom(
				"",
				u.Lobby.Id,
				eventMatch,
				matchEvent{
					Card: *card,
				},
			)
		}

		newCard := u.Card()
		conn.Emit(eventCard, cardEvent{Card: *newCard})
	})

	s.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})
}

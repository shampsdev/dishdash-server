package swipe

import (
	"context"
	"encoding/json"
	"log"

	"dishdash.ru/internal/dto"

	"dishdash.ru/internal/usecase"

	socketio "github.com/googollee/go-socket.io"
)

func SetupLobby(wsServer *socketio.Server, useCases usecase.Cases) {
	wsServer.OnConnect("", func(conn socketio.Conn) error {
		conn.SetContext("")
		log.Println("connected:", conn.ID())
		return nil
	})

	wsServer.OnEvent("", "echo", func(s socketio.Conn, msg string) {
		log.Println("echo:", msg)
		s.Emit("echo", msg)
	})

	wsServer.OnEvent("", eventJoinLobby, func(conn socketio.Conn, msg string) {
		var e joinLobbyEvent
		err := json.Unmarshal([]byte(msg), &e)
		if err != nil {
			_ = conn.Close()
			return
		}
		_, ok := conn.Context().(user)
		if ok {
			return
		}

		domLobby, err := useCases.Lobby.GetLobbyByID(context.Background(), e.LobbyID)
		if err != nil {
			_ = conn.Close()
			return
		}

		lobby, err := findLobby(domLobby, useCases.Card)
		if err != nil {
			_ = conn.Close()
			return
		}

		u := &user{
			ID:     conn.ID(),
			lobby:  lobby,
			swipes: nil,
			conn:   conn,
		}

		lobby.registerUser(u)
		conn.SetContext(u)

		firstCard := u.takeCard()
		conn.Emit(eventCard, cardEvent{Card: firstCard.ToDto()})
	})

	wsServer.OnEvent("", eventSwipe, func(conn socketio.Conn, msg string) {
		var swipeEvent swipeEvent
		err := json.Unmarshal([]byte(msg), &swipeEvent)
		if err != nil {
			conn.Close()
			return
		}

		u, ok := conn.Context().(*user)
		if !ok {
			conn.Close()
			return
		}

		swipe := swipe{
			T:    swipeEvent.SwipeType,
			Card: u.takeCard(),
		}
		if swipe.T == dto.LIKE {
			conn.Emit(eventMatch, matchEvent{
				Card: u.takeCard().ToDto(),
			})
		}
		u.swipe(swipe)

		newCard := u.takeCard()
		conn.Emit(eventCard, cardEvent{Card: newCard.ToDto()})
	})

	wsServer.OnDisconnect("", func(s socketio.Conn, reason string) {
		u, ok := s.Context().(*user)
		if ok {
			u.lobby.unregisterUser(u)
		}
		log.Println("disconnected:", s.ID(), "reason:", reason)
	})
}

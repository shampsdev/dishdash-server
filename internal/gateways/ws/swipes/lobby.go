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
		log.Println("received connection bro")
		var joinEvent joinLobbyEvent
		err := json.Unmarshal([]byte(msg), &joinEvent)
		if err != nil {
			_ = conn.Close()
			return
		}
		_, ok := conn.Context().(*entities.User)
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
			eventUserJoined,
			userJoinEvent{
				Name:   u.Name,
				Avatar: u.Avatar,
			},
		)
	})

	s.OnEvent("/", eventStartSwipes, func(conn socketio.Conn, msg string) {
		user, ok := conn.Context().(*entities.User)
		if !ok {
			log.Println("user not registered, disconnected")
			_ = conn.Close()
			return
		}

		s.ForEach("/", user.Lobby.Id, func(c socketio.Conn) {
			roomUser, ok := c.Context().(*entities.User)
			if !ok {
				log.Println("Failed to retrieve user from connection context.")
			}

			firstCard := roomUser.Card()
			c.Emit(eventCard, cardEvent{Card: *firstCard})
		})

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

		user.Lobby.UpdateSettings(domain.LobbySettings{
			PriceMin:    updateEvent.PriceMin,
			PriceMax:    updateEvent.PriceMax,
			MaxDistance: updateEvent.MaxDistance,
			Tags:        updateEvent.Tags,
		})

		s.BroadcastToRoom(
			"",
			user.Lobby.Id,
			eventSettingsUpdate,
			updateEvent,
		)
	})

	// s.OnEvent("", "getLobbyUpdate", func(conn socketio.Conn, msg string) {
	// 	u, ok := conn.Context().(*entities.User)
	// 	if !ok {
	// 		log.Println("user not registered, disconnected")
	// 		conn.Close()
	// 		return
	// 	}

	// 	// conn.Emit("lobbyUpdate", &lobbyUpdateEvent{

	// 	// })
	// })

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

		match := u.Swipe(swipeEvent.SwipeType)
		if match != nil {
			s.BroadcastToRoom(
				"",
				u.Lobby.Id,
				eventMatch,
				matchEvent{
					Id:   match.ID,
					Card: *card,
				},
			)

			vote := entities.NewVote(2, func(vote *entities.Vote, results []int) {
				sum := 0
				for _, number := range results {
					sum += number
				}
				if sum == len(u.Lobby.GetUsers()) {
					vote.FinalizeVote()
				}
			}, func(results []int) {
				s.BroadcastToRoom(
					"",
					u.Lobby.Id,
					eventRelaseMatch,
				)
			})

			u.Lobby.RegisterVote(vote, match.ID)
		}

		newCard := u.Card()
		conn.Emit(eventCard, cardEvent{Card: *newCard})
	})

	s.OnEvent("/", eventVote, func(conn socketio.Conn, msg string) {
		var voteEvent voteEvent
		err := json.Unmarshal([]byte(msg), &voteEvent)
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

		v := u.Lobby.GetVoteById(voteEvent.VoteId)

		v.Vote(int(voteEvent.VoteOption))
		log.Println("user", v)
	})

	s.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})
}

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

func SetupHandlers(sio *socketio.Server, cases usecase.Cases) {
	s := NewServer(sio)

	s.SIO.OnDisconnect("/", func(conn socketio.Conn, msg string) {
		c, ok := s.GetContext(conn)
		if !ok {
			log.Infof("disconnected with msg \"%s\" ", msg)
			return
		}

		defer func() {
			s.Metrics.ActiveConnections.Dec()
			c.Log.Info("Leave room")
			if c.Room != nil && c.Room.Empty() {
				err := cases.RoomRepo.DeleteRoom(context.Background(), c.Room.ID())
				if err != nil {
					c.Log.WithError(err).Error("error while deleting room")
				}
			}
		}()

		if c.Room == nil {
			log.Warn("room not found in context while disconnect")
			return
		}

		err := c.Room.RemoveUser(c.User.ID)
		if err != nil {
			c.HandleError(fmt.Errorf("error while removing user from room: %w", err))
		}

		s.SIO.LeaveRoom("", c.Room.ID(), conn)
		s.SIO.BroadcastToRoom("", c.Room.ID(), event.UserLeft,
			event.UserLeftEvent{
				ID:     c.User.ID,
				Name:   c.User.Name,
				Avatar: c.User.Avatar,
			})
	})

	s.On(event.JoinLobby, EventOpts{
		Allowed: []domain.LobbyState{domain.InLobby, domain.Finished},
	},
		func(c *Context, joinEvent event.JoinLobbyEvent) {
			c.Log = log.WithFields(log.Fields{
				"user":  joinEvent.UserID,
				"room":  joinEvent.LobbyID,
				"event": event.JoinLobby,
			})

			user, err := cases.User.GetUserByID(context.Background(), joinEvent.UserID)
			if err != nil {
				c.HandleError(fmt.Errorf("error while getting user: %w", err))
				return
			}

			c.User = user

			// room, err := cases.RoomRepo.GetRoom(context.Background(), joinEvent.LobbyID)
			var room *usecase.Room
			if err != nil {
				c.HandleError(fmt.Errorf("error while getting room: %w", err))
				return
			}

			c.Room = room

			settings := c.Room.Settings()
			c.Emit(event.SettingsUpdate, event.SettingsUpdateEvent{
				Location:    settings.Location,
				PriceMin:    settings.PriceAvg - 300,
				PriceMax:    settings.PriceAvg + 300,
				MaxDistance: 4000,
				Tags:        settings.Tags,
			})

			if room.Finished() {
				c.Emit(event.Finish, event.FinishEvent{
					Result:  c.Room.Result(),
					Matches: c.Room.Matches(),
				})
				return
			}

			for _, v := range room.Votes() {
				c.Emit(event.VoteAnnounce, v)
			}

			err = room.AddUser(user)
			if err != nil {
				c.HandleError(fmt.Errorf("error while adding user to room: %w", err))
				return
			}

			c.Conn.Join(room.ID())
			c.lock.Lock()
			c.Conn.SetContext(c)
			c.lock.Unlock()
			c.Log.Info("User joined")

			broadcastToOthersInRoom(s.SIO, c.User.ID, c.Room.ID(), event.UserJoined,
				event.UserJoinedEvent{
					ID:     user.ID,
					Name:   user.Name,
					Avatar: user.Avatar,
				})

			for _, u := range c.Room.Users() {
				c.Emit(event.UserJoined, event.UserJoinedEvent{
					ID:     u.ID,
					Name:   u.Name,
					Avatar: u.Avatar,
				})
			}

			if c.Room.Swiping() {
				c.Emit(event.StartSwipes)
			}

			if c.Room.Finished() {
				c.Emit(event.Finish, event.FinishEvent{
					Result:  c.Room.Result(),
					Matches: c.Room.Matches(),
				})
			}
		})

	s.On(event.SettingsUpdate, EventOpts{
		Allowed: []domain.LobbyState{domain.InLobby},
	},
		func(c *Context, se event.SettingsUpdateEvent) {
			ctx := context.Background()
			err := c.Room.UpdateLobbySettings(ctx, se.Location, (se.PriceMax+se.PriceMax)/2, se.Tags, nil, se.RecommendationOpts)
			if err != nil {
				c.HandleError(fmt.Errorf("error while updating lobby: %w", err))
				return
			}

			se.UserID = c.User.ID
			s.SIO.BroadcastToRoom("", c.Room.ID(), event.SettingsUpdate, se)
		})

	s.On(event.StartSwipes, EventOpts{
		Allowed: []domain.LobbyState{domain.InLobby},
	},
		func(c *Context) {
			err := c.Room.StartSwipes(context.Background())
			if err != nil {
				c.HandleError(fmt.Errorf("error while starting swipes: %w", err))
				return
			}

			c.Log.Info("Start swipes")
			s.SIO.BroadcastToRoom("", c.Room.ID(), event.StartSwipes)
			s.SIO.ForEach("", c.Room.ID(), func(conn socketio.Conn) {
				c, ok := s.GetContext(conn)
				if !ok {
					return
				}

				p := c.Room.GetNextPlaceForUser(c.User.ID)
				c.Emit(event.Place, event.PlaceEvent{
					ID:   p.ID,
					Card: p,
				})
			})
		})

	s.On(event.Swipe, EventOpts{
		Allowed: []domain.LobbyState{domain.Swiping},
	},
		func(c *Context, se event.SwipeEvent) {
			v, err := c.Room.Swipe(c.User.ID, c.Room.GetNextPlaceForUser(c.User.ID).ID, se.SwipeType)
			if err != nil {
				c.HandleError(fmt.Errorf("error while swiping: %w", err))
				return
			}

			if v != nil {
				s.SIO.BroadcastToRoom("/", c.Room.ID(), event.VoteAnnounce, v)
			}
			p := c.Room.GetNextPlaceForUser(c.User.ID)
			c.Emit(event.Place, event.PlaceEvent{
				ID:   p.ID,
				Card: p,
			})
		})

	voteLock := sync.RWMutex{}
	s.On(event.Vote, EventOpts{
		Allowed: []domain.LobbyState{domain.Swiping},
	},
		func(c *Context, ve event.VoteEvent) {
			voteLock.Lock()
			defer voteLock.Unlock()
			res, err := c.Room.Vote(c.User.ID, ve.VoteID, ve.OptionID)
			if err != nil {
				c.HandleError(fmt.Errorf("error while voting: %w", err))
				return
			}
			s.SIO.BroadcastToRoom("/", c.Room.ID(), event.Voted, event.VotedEvent{
				VoteID:   ve.VoteID,
				OptionID: ve.OptionID,
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

			if res != nil {
				s.SIO.BroadcastToRoom("/", c.Room.ID(), event.VoteResult, event.VoteResultEvent{
					VoteID:   res.VoteID,
					OptionID: res.OptionID,
				})
			}

			if c.Room.Finished() {
				s.SIO.BroadcastToRoom("/", c.Room.ID(), event.Finish, event.FinishEvent{
					Result:  c.Room.Result(),
					Matches: c.Room.Matches(),
				})
			}
		})

	s.On(event.LeaveLobby, EventOpts{},
		func(c *Context) {
			if err := c.Conn.Close(); err != nil {
				log.WithError(err).Error("error while closing connection")
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
			c.Emit(event, args...)
		}
	})
}

package tests

import (
	"testing"
	"time"

	"dishdash.ru/e2e/pg_test"

	"dishdash.ru/internal/gateways/ws/event"

	"dishdash.ru/internal/domain"
	socketio "github.com/googollee/go-socket.io"
	"github.com/stretchr/testify/assert"
)

func LobbySwipe(t *testing.T) *SocketIOSession {
	user1 := postUserWithID(t, &domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})
	user2 := postUserWithID(t, &domain.User{ID: "id2", Name: "user2", Avatar: "avatar2"})

	lobby := findLobby(t)

	sioCli1, err := socketio.NewClient(SIOHost, nil)
	assert.NoError(t, err)

	sioCli2, err := socketio.NewClient(SIOHost, nil)
	assert.NoError(t, err)

	sioSess := newSocketIOSession()
	sioSess.addUser(user1.Name)
	sioSess.addUser(user2.Name)

	sioCli1.OnEvent(event.UserJoined, sioSess.sioAddFunc(user1.Name, event.UserJoined))
	sioCli2.OnEvent(event.UserJoined, sioSess.sioAddFunc(user2.Name, event.UserJoined))
	sioCli1.OnEvent(event.Place, sioSess.sioAddFunc(user1.Name, event.Place))
	sioCli2.OnEvent(event.Place, sioSess.sioAddFunc(user2.Name, event.Place))
	sioCli1.OnEvent(event.Match, sioSess.sioAddFunc(user1.Name, event.Match))
	sioCli2.OnEvent(event.Match, sioSess.sioAddFunc(user2.Name, event.Match))

	assert.NoError(t, sioCli1.Connect())
	assert.NoError(t, sioCli2.Connect())

	sioSess.newStep("Joining lobby")
	sioCli1.Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})
	time.Sleep(waitTime)
	sioCli2.Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user2.ID,
	})
	sioCli1.Emit(event.SettingsUpdate, event.SettingsUpdateEvent{
		PriceMin:    300,
		PriceMax:    300,
		MaxDistance: 4000,
		Tags:        []int64{pg_test.Tags[3].ID},
	})
	time.Sleep(waitTime)

	sioSess.newStep("Start swipes")
	sioCli1.Emit(event.StartSwipes)
	time.Sleep(waitTime)

	sioSess.newStep("Swipe like and dislike")
	sioCli1.Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	sioCli2.Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.DISLIKE})
	time.Sleep(waitTime)

	sioSess.newStep("Swipe both likes")
	sioCli1.Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	sioCli2.Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	time.Sleep(waitTime)

	assert.NoError(t, sioCli1.Close())
	assert.NoError(t, sioCli2.Close())

	return sioSess
}

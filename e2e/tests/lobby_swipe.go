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

	cli1, err := socketio.NewClient(SIOHost, nil)
	assert.NoError(t, err)

	cli2, err := socketio.NewClient(SIOHost, nil)
	assert.NoError(t, err)

	sioSess := newSocketIOSession()
	sioSess.addUser(user1.Name)
	sioSess.addUser(user2.Name)

	listenEvent := func(eventName string) {
		cli1.OnEvent(eventName, sioSess.sioAddFunc(user1.Name, eventName))
		cli2.OnEvent(eventName, sioSess.sioAddFunc(user2.Name, eventName))
	}

	listenEvent(event.Error)
	listenEvent(event.UserJoined)
	listenEvent(event.StartSwipes)
	listenEvent(event.SettingsUpdate)
	listenEvent(event.Place)
	listenEvent(event.Match)

	assert.NoError(t, cli1.Connect())
	assert.NoError(t, cli2.Connect())

	cli1Emit := emitFuncWithLog(cli1, user1.Name)
	cli2Emit := emitFuncWithLog(cli2, user2.Name)

	sioSess.newStep("User1 join lobby")
	cli1Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})
	time.Sleep(waitTime)

	sioSess.newStep("Settings update")
	cli1Emit(event.SettingsUpdate, event.SettingsUpdateEvent{
		PriceMin:    300,
		PriceMax:    300,
		MaxDistance: 4000,
		Tags:        []int64{pg_test.Tags[3].ID},
	})
	time.Sleep(waitTime)

	sioSess.newStep("User2 join lobby")
	cli2Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user2.ID,
	})
	time.Sleep(waitTime)

	sioSess.newStep("Start swipes")
	cli1Emit(event.StartSwipes)
	time.Sleep(waitTime)

	sioSess.newStep("Swipe like and dislike")
	cli1Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	cli2Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.DISLIKE})
	time.Sleep(waitTime)

	sioSess.newStep("Swipe both likes")
	cli1Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	cli2Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	time.Sleep(waitTime)

	assert.NoError(t, cli1.Close())
	assert.NoError(t, cli2.Close())

	return sioSess
}

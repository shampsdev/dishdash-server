package tests

import (
	"testing"
	"time"

	"dishdash.ru/e2e/sdk"
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/gateways/ws/event"

	socketio "github.com/googollee/go-socket.io"
	"github.com/stretchr/testify/assert"
)

func Test_LobbySwipe(t *testing.T) {
	sdk.RunSessionTest(t, sdk.SessionTest{
		GoldenFile: "lobby_swipe",
		Run:        LobbySwipe,
	})
}

func LobbySwipe(t *testing.T) *sdk.SocketIOSession {
	user1 := sdk.PostUserWithID(t, &domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})
	user2 := sdk.PostUserWithID(t, &domain.User{ID: "id2", Name: "user2", Avatar: "avatar2"})

	lobby := sdk.FindLobby(t)

	cli1, err := socketio.NewClient(sdk.SIOHost, nil)
	assert.NoError(t, err)

	cli2, err := socketio.NewClient(sdk.SIOHost, nil)
	assert.NoError(t, err)

	sioSess := sdk.NewSocketIOSession()
	sioSess.AddUser(user1.Name)
	sioSess.AddUser(user2.Name)

	listenEvent := func(eventName string) {
		cli1.OnEvent(eventName, sioSess.SioAddFunc(user1.Name, eventName))
		cli2.OnEvent(eventName, sioSess.SioAddFunc(user2.Name, eventName))
	}

	listenEvent(event.Error)
	listenEvent(event.UserJoined)
	listenEvent(event.StartSwipes)
	listenEvent(event.SettingsUpdate)
	listenEvent(event.Place)
	listenEvent(event.Match)

	assert.NoError(t, cli1.Connect())
	assert.NoError(t, cli2.Connect())

	cli1Emit := sdk.EmitWithLogFunc(cli1, user1.Name)
	cli2Emit := sdk.EmitWithLogFunc(cli2, user2.Name)

	sioSess.NewStep("User1 join lobby")
	cli1Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})
	time.Sleep(sdk.WaitTime)

	sioSess.NewStep("Settings update")
	cli1Emit(event.SettingsUpdate, event.SettingsUpdateEvent{
		PriceMin:    300,
		PriceMax:    300,
		MaxDistance: 4000,
		Tags:        []int64{4},
	})
	time.Sleep(sdk.WaitTime)

	sioSess.NewStep("User2 join lobby")
	cli2Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user2.ID,
	})
	time.Sleep(sdk.WaitTime)

	sioSess.NewStep("Start swipes")
	cli1Emit(event.StartSwipes)
	time.Sleep(sdk.WaitTime)

	sioSess.NewStep("Swipe like and dislike")
	cli1Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	cli2Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.DISLIKE})
	time.Sleep(sdk.WaitTime)

	sioSess.NewStep("Swipe both likes")
	cli1Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	cli2Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	time.Sleep(sdk.WaitTime)

	assert.NoError(t, cli1.Close())
	assert.NoError(t, cli2.Close())

	return sioSess
}

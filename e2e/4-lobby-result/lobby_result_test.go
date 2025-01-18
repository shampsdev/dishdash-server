package tests

import (
	"testing"

	"dishdash.ru/e2e/sdk"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase/event"
	socketio "github.com/googollee/go-socket.io"
	"github.com/stretchr/testify/assert"
)

func Test_LobbyResult(t *testing.T) {
	sdk.RunSessionTest(t, sdk.SessionTest{
		GoldenFile: "lobby_result",
		Run:        LobbyResult,
	})
}

func LobbyResult(t *testing.T) *sdk.SocketIOSession {
	user1 := sdk.PostUserWithID(t, &domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})
	user2 := sdk.PostUserWithID(t, &domain.User{ID: "id2", Name: "user2", Avatar: "avatar2"})

	lobby := sdk.FindLobby(t)
	sdk.Sleep()

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
	listenEvent(event.VoteAnnounceEvent)
	listenEvent(event.SettingsUpdateEvent)

	assert.NoError(t, cli1.Connect())
	assert.NoError(t, cli2.Connect())

	cli1Emit := sdk.EmitWithLogFunc(cli1, user1.Name)
	cli2Emit := sdk.EmitWithLogFunc(cli2, user2.Name)

	sioSess.NewStep("Joining lobby")
	cli1Emit(event.JoinLobbyEvent, event.JoinLobby{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})
	cli2Emit(event.JoinLobbyEvent, event.JoinLobby{
		LobbyID: lobby.ID,
		UserID:  user2.ID,
	})
	cli1Emit(event.SettingsUpdateEvent, event.SettingsUpdate{
		Location:    lobby.Location,
		PriceMin:    300,
		PriceMax:    300,
		MaxDistance: 4000,
		Tags:        []int64{4},
	})
	sdk.Sleep()

	sioSess.NewStep("Start swipes")
	cli1Emit(event.StartSwipesEvent)
	sdk.Sleep()

	sioSess.NewStep("Swipe like and dislike")
	cli1Emit(event.SwipeEvent, event.Swipe{SwipeType: domain.DISLIKE})
	cli2Emit(event.SwipeEvent, event.Swipe{SwipeType: domain.LIKE})
	sdk.Sleep()

	sioSess.NewStep("Swipe both likes")
	cli1Emit(event.SwipeEvent, event.Swipe{SwipeType: domain.LIKE})
	cli2Emit(event.SwipeEvent, event.Swipe{SwipeType: domain.LIKE})
	sdk.Sleep()

	cli1Emit(event.LeaveLobbyEvent)
	cli2Emit(event.LeaveLobbyEvent)
	sdk.Sleep()

	// leave to check if lobby will be finished
	assert.NoError(t, cli1.Close())
	assert.NoError(t, cli2.Close())

	cli3, err := socketio.NewClient(sdk.SIOHost, nil)
	assert.NoError(t, err)

	listenEvent = func(eventName string) {
		cli3.OnEvent(eventName, sioSess.SioAddFunc(user1.Name, eventName))
	}
	listenEvent(event.JoinLobbyEvent)
	listenEvent(event.FinishEvent)
	listenEvent(event.SettingsUpdateEvent)
	listenEvent(event.ErrorEvent)

	sioSess.NewStep("Rejoin")
	assert.NoError(t, cli3.Connect())
	cli3Emit := sdk.EmitWithLogFunc(cli3, user1.Name)
	cli3Emit(event.JoinLobbyEvent, event.JoinLobby{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})

	sdk.Sleep()
	assert.NoError(t, cli3.Close())

	return sioSess
}

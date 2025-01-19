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
		GoldenFile: "lobby_rejoin",
		Run:        LobbyResult,
	})
}

func LobbyResult(t *testing.T) *sdk.SocketIOSession {
	user1 := sdk.PostUserWithID(t, &domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})

	lobby := sdk.FindLobby(t)
	sdk.Sleep()

	cli1, err := socketio.NewClient(sdk.SIOHost, nil)
	assert.NoError(t, err)

	sioSess := sdk.NewSocketIOSession()
	sioSess.AddUser(user1.Name)

	listenEvent := func(eventName string) {
		cli1.OnEvent(eventName, sioSess.SioAddFunc(user1.Name, eventName))
	}
	listenEvent(event.VoteAnnounceEvent)
	listenEvent(event.SettingsUpdateEvent)

	assert.NoError(t, cli1.Connect())

	cli1Emit := sdk.EmitWithLogFunc(cli1, user1.Name)

	sioSess.NewStep("Joining lobby")
	cli1Emit(event.JoinLobbyEvent, event.JoinLobby{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
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

	sioSess.NewStep("Swipe dislike")
	cli1Emit(event.SwipeEvent, event.Swipe{SwipeType: domain.DISLIKE})
	sdk.Sleep()

	sioSess.NewStep("Swipe like")
	cli1Emit(event.SwipeEvent, event.Swipe{SwipeType: domain.LIKE})
	sdk.Sleep()

	cli1Emit(event.LeaveLobbyEvent)
	sdk.Sleep()

	// leave to check if lobby will be finished
	assert.NoError(t, cli1.Close())

	cli1, err = socketio.NewClient(sdk.SIOHost, nil)
	assert.NoError(t, err)

	listenEvent = func(eventName string) {
		cli1.OnEvent(eventName, sioSess.SioAddFunc(user1.Name, eventName))
	}
	listenEvent(event.JoinLobbyEvent)
	listenEvent(event.PlaceEvent)
	listenEvent(event.VoteAnnounceEvent)
	listenEvent(event.SettingsUpdateEvent)
	listenEvent(event.ErrorEvent)

	sioSess.NewStep("Rejoin")
	assert.NoError(t, cli1.Connect())
	cli3Emit := sdk.EmitWithLogFunc(cli1, user1.Name)
	cli3Emit(event.JoinLobbyEvent, event.JoinLobby{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})

	sdk.Sleep()
	assert.NoError(t, cli1.Close())

	return sioSess
}

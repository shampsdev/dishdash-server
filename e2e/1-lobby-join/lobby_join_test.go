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

func Test_LobbyJoin(t *testing.T) {
	sdk.RunSessionTest(t, sdk.SessionTest{
		GoldenFile: "lobby_join",
		Run:        LobbyJoin,
	})
}

func LobbyJoin(t *testing.T) *sdk.SocketIOSession {
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

	cli1.OnEvent(event.UserJoined, sioSess.SioAddFunc(user1.Name, event.UserJoined))
	cli2.OnEvent(event.UserJoined, sioSess.SioAddFunc(user2.Name, event.UserJoined))

	assert.NoError(t, cli1.Connect())
	assert.NoError(t, cli2.Connect())

	sioSess.NewStep("Joining lobby")
	cli1Emit := sdk.EmitWithLogFunc(cli1, user1.Name)
	cli2Emit := sdk.EmitWithLogFunc(cli2, user2.Name)

	cli1Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})
	time.Sleep(sdk.WaitTime)
	cli2Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user2.ID,
	})
	time.Sleep(sdk.WaitTime)

	assert.NoError(t, cli1.Close())
	assert.NoError(t, cli2.Close())

	return sioSess
}

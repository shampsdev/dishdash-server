package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/gateways/ws/event"
	"dishdash.ru/internal/usecase"

	socketio "github.com/googollee/go-socket.io"
	"github.com/stretchr/testify/assert"
)

func LobbyJoin(t *testing.T) *SocketIOSession {
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

	cli1.OnEvent(event.UserJoined, sioSess.sioAddFunc(user1.Name, event.UserJoined))
	cli2.OnEvent(event.UserJoined, sioSess.sioAddFunc(user2.Name, event.UserJoined))

	assert.NoError(t, cli1.Connect())
	assert.NoError(t, cli2.Connect())

	sioSess.newStep("Joining lobby")
	cli1Emit := emitWithLogFunc(cli1, user1.Name)
	cli2Emit := emitWithLogFunc(cli2, user2.Name)

	cli1Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})
	time.Sleep(waitTime)
	cli2Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user2.ID,
	})
	time.Sleep(waitTime)

	assert.NoError(t, cli1.Close())
	assert.NoError(t, cli2.Close())

	return sioSess
}

func findLobby(t *testing.T) *domain.Lobby {
	findLobbyInput := usecase.FindLobbyInput{
		Dist: 0,
		// ИТМО - Кронверкский проспект, 49
		Location: domain.Coordinate{Lon: 30.310011, Lat: 59.956363},
	}
	b, err := json.Marshal(findLobbyInput)
	assert.NoError(t, err)

	resp, err := httpClient.Post(fmt.Sprintf("%s/lobbies/find", ApiHost), "application/json", bytes.NewReader(b))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	lobby := &domain.Lobby{}
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(lobby))
	return lobby
}

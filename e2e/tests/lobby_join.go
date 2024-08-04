package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/gateways/ws/event"
	"dishdash.ru/internal/usecase"

	socketio "github.com/googollee/go-socket.io"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func JoinLobby(t *testing.T) {
	user1 := postUser(t, &domain.User{Name: "user1", Avatar: "avatar1"})
	user2 := postUser(t, &domain.User{Name: "user2", Avatar: "avatar2"})

	lobby := findLobby(t)

	sioCli1, err := socketio.NewClient(SIOHost, nil)
	assert.NoError(t, err)

	sioCli2, err := socketio.NewClient(SIOHost, nil)
	assert.NoError(t, err)

	e1Mu := &sync.Mutex{}
	var events1 []map[string]interface{}
	e2Mu := &sync.Mutex{}
	var events2 []map[string]interface{}

	sioCli1.OnEvent(event.UserJoined, func(_ socketio.Conn, data map[string]interface{}) {
		e1Mu.Lock()
		defer e1Mu.Unlock()
		events1 = append(events1, data)
	})
	sioCli2.OnEvent(event.UserJoined, func(_ socketio.Conn, data map[string]interface{}) {
		e2Mu.Lock()
		defer e2Mu.Unlock()
		events2 = append(events2, data)
	})

	assert.NoError(t, sioCli1.Connect())
	assert.NoError(t, sioCli2.Connect())

	sioCli1.Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})
	sioCli2.Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user2.ID,
	})
	time.Sleep(5 * time.Second)

	e1Mu.Lock()
	assert.Equal(t, 1, len(events1))
	userJoinedEvent := event.UserJoinedEvent{}
	err = mapstructure.Decode(events1[0], &userJoinedEvent)
	assert.NoError(t, err)
	assert.Equal(t, event.UserJoinedEvent{
		ID:     user2.ID,
		Name:   user2.Name,
		Avatar: user2.Avatar,
	}, userJoinedEvent)
	e1Mu.Unlock()

	assert.NoError(t, sioCli1.Close())
	assert.NoError(t, sioCli2.Close())
}

func findLobby(t *testing.T) *domain.Lobby {
	findLobbyInput := usecase.FindLobbyInput{
		Dist:     0,
		Location: domain.Coordinate{Lon: 2, Lat: 2},
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

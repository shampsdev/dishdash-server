package tests

import (
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/gateways/ws/event"
	socketio "github.com/googollee/go-socket.io"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func SwipeLobby(t *testing.T) {
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

	f1 := func(_ socketio.Conn, data map[string]interface{}) {
		e1Mu.Lock()
		defer e1Mu.Unlock()
		events1 = append(events1, data)
	}
	f2 := func(_ socketio.Conn, data map[string]interface{}) {
		e2Mu.Lock()
		defer e2Mu.Unlock()
		events2 = append(events2, data)
	}

	sioCli1.OnEvent(event.UserJoined, f1)
	sioCli2.OnEvent(event.UserJoined, f2)
	sioCli1.OnEvent(event.Place, f1)
	sioCli2.OnEvent(event.Place, f2)
	sioCli1.OnEvent(event.Match, f1)
	sioCli2.OnEvent(event.Match, f2)

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

	sioCli1.Emit(event.StartSwipes)
	time.Sleep(10 * time.Second)

	e1Mu.Lock()
	assert.Equal(t, 2, len(events1))
	var pe1 event.PlaceEvent
	assert.NoError(t, mapStructureDecode(events1[1], &pe1))
	assert.NotEmpty(t, pe1.Card.Title)
	e1Mu.Unlock()

	sioCli1.Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	sioCli2.Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.DISLIKE})
	time.Sleep(5 * time.Second)

	e1Mu.Lock()
	assert.Equal(t, 3, len(events1))
	var pe2 event.PlaceEvent
	assert.NoError(t, mapStructureDecode(events1[2], &pe2))
	assert.NotEmpty(t, pe2.Card.Title)
	e1Mu.Unlock()

	sioCli1.Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	sioCli2.Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	time.Sleep(5 * time.Second)

	e1Mu.Lock()
	assert.Equal(t, 4, len(events1))
	var me event.MatchEvent
	assert.NoError(t, mapStructureDecode(events1[3], &me))
	assert.Equal(t, *me.Card, *pe2.Card)
	e1Mu.Unlock()

	assert.NoError(t, sioCli1.Close())
	assert.NoError(t, sioCli2.Close())
}

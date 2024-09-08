package tests

import (
	"testing"
	"time"

	"dishdash.ru/e2e/pg_test"
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/gateways/ws/event"
	"dishdash.ru/internal/usecase"
	socketio "github.com/googollee/go-socket.io"
	"github.com/stretchr/testify/assert"
)

func DiffTimeJoin(t *testing.T) *SocketIOSession {
	user1 := postUserWithID(t, &domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})
	user2 := postUserWithID(t, &domain.User{ID: "id2", Name: "user2", Avatar: "avatar2"})
	user3 := postUserWithID(t, &domain.User{ID: "id3", Name: "user3", Avatar: "avatar3"})

	lobby := findLobby(t)

	sioSess := newSocketIOSession()
	sioSess.addUser(user1.Name)
	sioSess.addUser(user2.Name)
	sioSess.addUser(user3.Name)

	cli1, err := socketio.NewClient(SIOHost, nil)
	assert.NoError(t, err)
	cli1Emit := emitWithLogFunc(cli1, user1.Name)

	cli2, err := socketio.NewClient(SIOHost, nil)
	assert.NoError(t, err)
	cli2Emit := emitWithLogFunc(cli2, user2.Name)

	cli3, err := socketio.NewClient(SIOHost, nil)
	assert.NoError(t, err)
	cli3Emit := emitWithLogFunc(cli3, user3.Name)

	listenEvent := func(eventName string) {
		cli1.OnEvent(eventName, sioSess.sioAddFunc(user1.Name, eventName))
		cli2.OnEvent(eventName, sioSess.sioAddFunc(user2.Name, eventName))
		cli3.OnEvent(eventName, sioSess.sioAddFunc(user3.Name, eventName))
	}

	listenEvent(event.Error)
	listenEvent(event.UserJoined)
	listenEvent(event.StartSwipes)
	listenEvent(event.SettingsUpdate)
	listenEvent(event.Finish)

	sioSess.newStep("User 1 joins lobby")
	assert.NoError(t, cli1.Connect())
	cli1Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})
	time.Sleep(waitTime)
	cli1Emit(event.SettingsUpdate, event.SettingsUpdateEvent{
		PriceMin:    300,
		PriceMax:    300,
		MaxDistance: 4000,
		Tags:        []int64{pg_test.Tags[3].ID},
	})
	time.Sleep(waitTime)

	sioSess.newStep("Start swipes")
	cli1Emit(event.StartSwipes)
	time.Sleep(waitTime)

	sioSess.newStep("User 2 joins lobby")
	assert.NoError(t, cli2.Connect())

	cli2Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user2.ID,
	})
	time.Sleep(waitTime)

	sioSess.newStep("Swipe both like")
	cli1Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	cli2Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	time.Sleep(waitTime)

	sioSess.newStep("Both accepted match")
	cli1Emit(event.Vote, event.VoteEvent{
		ID:     0,
		Option: usecase.VoteLike,
	})
	cli2Emit(event.Vote, event.VoteEvent{
		ID:     0,
		Option: usecase.VoteLike,
	})
	time.Sleep(waitTime)

	sioSess.newStep("User 3 joins lobby")
	assert.NoError(t, cli3.Connect())
	cli3Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user3.ID,
	})
	time.Sleep(waitTime)

	return sioSess
}

package tests

import (
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/gateways/ws/event"
	"dishdash.ru/internal/usecase"
	socketio "github.com/googollee/go-socket.io"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func LobbyVote(t *testing.T) *SocketIOSession {
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

	sioCli1.OnEvent(event.Match, sioSess.sioAddFunc(user1.Name, event.Match))
	sioCli2.OnEvent(event.Match, sioSess.sioAddFunc(user2.Name, event.Match))
	sioCli1.OnEvent(event.Voted, sioSess.sioAddFunc(user1.Name, event.Voted))
	sioCli2.OnEvent(event.Voted, sioSess.sioAddFunc(user2.Name, event.Voted))
	sioCli1.OnEvent(event.ReleaseMatch, sioSess.sioAddFunc(user1.Name, event.ReleaseMatch))
	sioCli2.OnEvent(event.ReleaseMatch, sioSess.sioAddFunc(user2.Name, event.ReleaseMatch))
	sioCli1.OnEvent(event.Finish, sioSess.sioAddFunc(user1.Name, event.Finish))
	sioCli2.OnEvent(event.Finish, sioSess.sioAddFunc(user2.Name, event.Finish))

	assert.NoError(t, sioCli1.Connect())
	assert.NoError(t, sioCli2.Connect())

	sioSess.newStep("Joining lobby")
	sioCli1.Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})
	sioCli2.Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user2.ID,
	})
	time.Sleep(waitTime)

	sioSess.newStep("Start swipes")
	sioCli1.Emit(event.StartSwipes)
	time.Sleep(waitTime)

	sioSess.newStep("Swipe both likes (1)")
	sioCli1.Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	sioCli2.Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	time.Sleep(waitTime)

	sioSess.newStep("Vote like and dislike")
	sioCli1.Emit(event.Vote, event.VoteEvent{ID: 0, Option: usecase.VoteLike})
	sioCli2.Emit(event.Vote, event.VoteEvent{ID: 0, Option: usecase.VoteDislike})
	time.Sleep(waitTime)

	sioSess.newStep("Swipe both likes (2)")
	sioCli1.Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	sioCli2.Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	time.Sleep(waitTime)

	sioSess.newStep("Vote both likes")
	sioCli1.Emit(event.Vote, event.VoteEvent{ID: 1, Option: usecase.VoteLike})
	sioCli2.Emit(event.Vote, event.VoteEvent{ID: 1, Option: usecase.VoteLike})
	time.Sleep(waitTime)

	assert.NoError(t, sioCli1.Close())
	assert.NoError(t, sioCli2.Close())

	return sioSess
}

package tests

import (
	"testing"

	"dishdash.ru/e2e/sdk"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/gateways/ws/event"
	"dishdash.ru/internal/usecase"
	socketio "github.com/googollee/go-socket.io"
	"github.com/stretchr/testify/assert"
)

func Test_LobbyVote(t *testing.T) {
	sdk.RunSessionTest(t, sdk.SessionTest{
		GoldenFile: "lobby_vote",
		Run:        LobbyVote,
	})
}

func LobbyVote(t *testing.T) *sdk.SocketIOSession {
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

	listenEvent(event.VoteAnnounce)
	listenEvent(event.Voted)
	listenEvent(event.VoteResult)
	listenEvent(event.Finish)

	assert.NoError(t, cli1.Connect())
	assert.NoError(t, cli2.Connect())

	cli1Emit := sdk.EmitWithLogFunc(cli1, user1.Name)
	cli2Emit := sdk.EmitWithLogFunc(cli2, user2.Name)

	sioSess.NewStep("Joining lobby")
	cli1Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user1.ID,
	})
	cli2Emit(event.JoinLobby, event.JoinLobbyEvent{
		LobbyID: lobby.ID,
		UserID:  user2.ID,
	})
	cli1Emit(event.SettingsUpdate, event.SettingsUpdateEvent{
		Location:    lobby.Location,
		PriceMin:    300,
		PriceMax:    300,
		MaxDistance: 4000,
		Tags:        []int64{4},
	})
	sdk.Sleep()

	sioSess.NewStep("Start swipes")
	cli1Emit(event.StartSwipes)
	sdk.Sleep()

	sioSess.NewStep("Swipe both likes (1)")
	cli1Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	cli2Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	sdk.Sleep()

	sioSess.NewStep("Vote like and dislike")
	cli1Emit(event.Vote, event.VoteEvent{VoteID: 1, OptionID: usecase.OptionIDLike})
	cli2Emit(event.Vote, event.VoteEvent{VoteID: 1, OptionID: usecase.OptionIDDislike})
	sdk.Sleep()

	sioSess.NewStep("Swipe both likes (2)")
	cli1Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	cli2Emit(event.Swipe, event.SwipeEvent{SwipeType: domain.LIKE})
	sdk.Sleep()

	sioSess.NewStep("Vote both likes")
	cli1Emit(event.Vote, event.VoteEvent{VoteID: 2, OptionID: usecase.OptionIDLike})
	cli2Emit(event.Vote, event.VoteEvent{VoteID: 2, OptionID: usecase.OptionIDLike})
	sdk.Sleep()

	assert.NoError(t, cli1.Close())
	assert.NoError(t, cli2.Close())

	return sioSess
}

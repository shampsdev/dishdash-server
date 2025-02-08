package tests

import (
	"testing"

	"dishdash.ru/e2e/framework"
	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/usecase/event"

	"github.com/stretchr/testify/assert"
)

var fw = framework.MustInit()

func TestMain(m *testing.M) {
	fw.RecordEvents(
		event.ErrorEvent,
		event.VoteAnnounceEvent,
		event.VotedEvent,
		event.VoteResultEvent,
		event.FinishEvent,
	)
	fw.TestMain(m)
}

func Test(t *testing.T) {
	cli1 := fw.MustNewClient(&domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})
	cli2 := fw.MustNewClient(&domain.User{ID: "id2", Name: "user2", Avatar: "avatar2"})
	lobby := fw.MustFindLobby()

	fw.Step("Joining lobby", func() {
		cli1.JoinLobby(lobby)
		cli2.JoinLobby(lobby)
	}, 8)

	fw.Step("Settings update", func() {
		cli1.Emit(event.SettingsUpdate{
			Location:    lobby.Location,
			PriceMin:    300,
			PriceMax:    300,
			MaxDistance: 4000,
			Tags:        []int64{4},
		})
	}, 2)

	fw.Step("Start swipes", func() {
		cli1.Emit(event.StartSwipes{})
	}, 4)

	fw.Step("Swipe both likes (1)", func() {
		cli1.Emit(event.Swipe{SwipeType: domain.LIKE})
		cli2.Emit(event.Swipe{SwipeType: domain.LIKE})
	}, 4)

	fw.Step("Vote like and dislike", func() {
		cli1.Emit(event.Vote{VoteID: 1, OptionID: event.OptionIDLike})
		cli2.Emit(event.Vote{VoteID: 1, OptionID: event.OptionIDDislike})
	}, 6)

	fw.Step("Swipe both likes (2)", func() {
		cli1.Emit(event.Swipe{SwipeType: domain.LIKE})
		cli2.Emit(event.Swipe{SwipeType: domain.LIKE})
	}, 4)

	fw.Step("Vote both likes", func() {
		cli1.Emit(event.Vote{VoteID: 2, OptionID: event.OptionIDLike})
		cli2.Emit(event.Vote{VoteID: 2, OptionID: event.OptionIDLike})
	}, 8)

	assert.NoError(t, cli1.Close())
	assert.NoError(t, cli2.Close())

	fw.AssertSession(t)
}

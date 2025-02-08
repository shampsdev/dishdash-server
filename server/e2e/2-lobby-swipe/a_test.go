package tests

import (
	"testing"

	"dishdash.ru/e2e/framework"
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase/event"

	"github.com/stretchr/testify/assert"
)

var fw = framework.MustInit()

func TestMain(m *testing.M) {
	fw.RecordEvents(
		event.ErrorEvent,
		event.UserJoinedEvent,
		event.SettingsUpdateEvent,
		event.StartSwipesEvent,
		event.PlaceEvent,
		event.VoteAnnounceEvent,
	)
	fw.TestMain(m)
}

func Test(t *testing.T) {
	cli1 := fw.MustNewClient(&domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})
	cli2 := fw.MustNewClient(&domain.User{ID: "id2", Name: "user2", Avatar: "avatar2"})
	lobby := fw.MustFindLobby()

	fw.Step("User1 join lobby", func() {
		cli1.JoinLobby(lobby)
	}, 3)

	fw.Step("Settings update", func() {
		cli1.Emit(event.SettingsUpdate{
			Location:    lobby.Location,
			PriceMin:    300,
			PriceMax:    300,
			MaxDistance: 4000,
			Tags:        []int64{4},
		})
	}, 1)

	fw.Step("User2 join lobby", func() {
		cli2.JoinLobby(lobby)
	}, 5)

	fw.Step("Start swipes", func() {
		cli1.Emit(event.StartSwipes{})
	}, 4)

	fw.Step("Swipe like and dislike", func() {
		cli1.Emit(event.Swipe{SwipeType: domain.LIKE})
		cli2.Emit(event.Swipe{SwipeType: domain.DISLIKE})
	}, 2)

	fw.Step("Swipe both likes", func() {
		cli1.Emit(event.Swipe{SwipeType: domain.LIKE})
		cli2.Emit(event.Swipe{SwipeType: domain.LIKE})
	}, 4)

	assert.NoError(t, cli1.Close())
	assert.NoError(t, cli2.Close())

	fw.AssertSession(t)
}

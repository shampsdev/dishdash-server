package tests

import (
	"testing"
	"time"

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
		event.VoteAnnounceEvent,
		event.PlaceEvent,
	)
	fw.TestMain(m)
}

func Test(t *testing.T) {
	cli1 := fw.MustNewClient(&domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})
	lobby := fw.MustFindLobby()

	fw.Step("Joining lobby", func() {
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

	fw.Step("Start swipes", func() {
		cli1.Emit(event.StartSwipes{})
	}, 2)

	fw.Step("Swipe dislike", func() {
		cli1.Emit(event.Swipe{SwipeType: domain.DISLIKE})
	}, 1)

	fw.Step("Swipe like", func() {
		cli1.Emit(event.Swipe{SwipeType: domain.LIKE})
	}, 2)

	cli1.Emit(event.LeaveLobby{})
	assert.NoError(t, cli1.Close())
	time.Sleep(2 * time.Second)

	cli1 = fw.MustNewClient(&domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})
	fw.Step("Rejoin", func() {
		cli1.JoinLobby(lobby)
	}, 5)

	assert.NoError(t, cli1.Close())

	fw.AssertSession(t)
}

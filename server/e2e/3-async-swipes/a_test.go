package tests

import (
	"testing"
	"time"

	"dishdash.ru/e2e/framework"
	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/usecase/event"

	"github.com/stretchr/testify/assert"
)

var fw = framework.MustInit()

func TestMain(m *testing.M) {
	fw.RecordEvents(
		event.ErrorEvent,
		event.UserJoinedEvent,
		event.SettingsUpdateEvent,
		event.StartSwipesEvent,
		event.CardsEvent,
		event.ResultsEvent,
		event.MatchEvent,
	)
	fw.UseShortener(event.CardsEvent, framework.CardsShortener)
	fw.UseShortener(event.ResultsEvent, framework.ResultsShortener)
	fw.TestMain(m)
}

func Test(t *testing.T) {
	cli1 := fw.MustNewClient(&domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})
	cli2 := fw.MustNewClient(&domain.User{ID: "id2", Name: "user2", Avatar: "avatar2"})
	lobby := fw.MustCreateLobby()

	fw.Step("User1 join lobby", func() {
		cli1.JoinLobby(lobby)
	}, 2)

	fw.Step("Settings update", func() {
		cli1.Emit(event.SettingsUpdate{
			Type: domain.ClassicPlacesLobbyType,
			ClassicPlaces: &domain.ClassicPlacesSettings{
				Location: lobby.Settings.ClassicPlaces.Location,
				PriceAvg: 300,
				Tags:     []int64{4},
			},
		})
	}, 1)

	fw.Step("Start swipes", func() {
		cli1.Emit(event.StartSwipes{})
	}, 2)

	fw.Step("User1 swipes", func() {
		cli1.Emit(event.Swipe{SwipeType: domain.LIKE})
		cli1.Emit(event.Swipe{SwipeType: domain.DISLIKE})
		cli1.Emit(event.Swipe{SwipeType: domain.LIKE})
	}, 7)

	fw.Step("User1 leave lobby", func() {
		cli1.Emit(event.LeaveLobby{})
	}, 0)
	time.Sleep(1 * time.Second)
	assert.NoError(t, cli1.Close())
	time.Sleep(1 * time.Second)

	fw.Step("User2 join lobby", func() {
		cli2.JoinLobby(lobby)
	}, 5)

	fw.Step("User2 swipes", func() {
		cli2.Emit(event.Swipe{SwipeType: domain.DISLIKE})
		cli2.Emit(event.Swipe{SwipeType: domain.LIKE})
		cli2.Emit(event.Swipe{SwipeType: domain.LIKE})
	}, 6)

	assert.NoError(t, cli2.Close())

	fw.AssertSession(t)
}

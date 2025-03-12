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

	fw.Step("User2 join lobby", func() {
		cli2.JoinLobby(lobby)
	}, 4)

	fw.Step("Start swipes", func() {
		cli1.Emit(event.StartSwipes{})
	}, 4)

	fw.Step("User1 swipe like (1)", func() {
		cli1.Emit(event.Swipe{SwipeType: domain.LIKE})
	}, 3)

	fw.Step("User2 swipe dislike (1)", func() {
		cli2.Emit(event.Swipe{SwipeType: domain.DISLIKE})
	}, 1)

	fw.Step("User1 swipe like (2)", func() {
		cli1.Emit(event.Swipe{SwipeType: domain.LIKE})
	}, 3)

	fw.Step("User2 swipe like (2)", func() {
		cli2.Emit(event.Swipe{SwipeType: domain.LIKE})
	}, 5)

	assert.NoError(t, cli1.Close())
	assert.NoError(t, cli2.Close())

	fw.AssertSession(t)
}

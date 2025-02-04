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
	lobby := fw.MustCreateLobby()

	fw.Step("Joining lobby", func() {
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

	fw.Step("Swipe dislike", func() {
		cli1.Emit(event.Swipe{SwipeType: domain.DISLIKE})
	}, 1)

	fw.Step("Swipe like", func() {
		cli1.Emit(event.Swipe{SwipeType: domain.LIKE})
	}, 3)

	cli1.Emit(event.LeaveLobby{})
	time.Sleep(1 * time.Second)
	assert.NoError(t, cli1.Close())
	time.Sleep(2 * time.Second)

	cli1 = fw.MustNewClient(&domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})
	fw.Step("Rejoin", func() {
		cli1.JoinLobby(lobby)
	}, 4)

	assert.NoError(t, cli1.Close())

	fw.AssertSession(t)
}

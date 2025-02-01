package framework

import (
	"os"
	"strconv"
	"testing"

	"dishdash.ru/e2e/framework/session"
	"dishdash.ru/internal/usecase/event"
	"github.com/stretchr/testify/assert"
)

var allEvents = map[string]struct{}{
	event.ErrorEvent: {},

	event.JoinLobbyEvent:      {},
	event.LeaveLobbyEvent:     {},
	event.UserJoinedEvent:     {},
	event.UserLeftEvent:       {},
	event.SettingsUpdateEvent: {},
	event.StartSwipesEvent:    {},
	event.PlaceEvent:          {},
	event.FinishEvent:         {},

	event.VoteAnnounceEvent: {},
	event.VoteEvent:         {},
	event.VotedEvent:        {},
	event.VoteResultEvent:   {},
}

func (fw *Framework) TestMain(m *testing.M) {
	if t, err := strconv.ParseBool(os.Getenv("E2E_TESTING")); err != nil || !t {
		fw.Log.Debugf("Skipping, because E2E_TESTING is not set to true")
		return
	}

	err := fw.SetupDB()
	if err != nil {
		panic(err)
	}

	m.Run()
}

func (fw *Framework) AssertSession(t *testing.T) {
	if update, err := strconv.ParseBool(os.Getenv("E2E_UPDATE_GOLDEN")); err == nil && update {
		assert.NoError(t, fw.Session.SaveToFile("golden.json"))
	} else {
		exp, err := session.LoadFromFile("golden.json")
		assert.NoError(t, err)
		session.AssertEqual(t, fw.Session, exp)
		if t.Failed() {
			assert.NoError(t, fw.Session.SaveToFile("ERROR_golden.json"))
		}
	}
}

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
	)
	fw.TestMain(m)
}

func Test(t *testing.T) {
	cli1 := fw.MustNewClient(&domain.User{ID: "id1", Name: "user1", Avatar: "avatar1"})
	cli2 := fw.MustNewClient(&domain.User{ID: "id2", Name: "user2", Avatar: "avatar2"})
	lobby := fw.MustCreateLobby()
	fw.Step("Joining lobby", func() {
		cli1.JoinLobby(lobby)
		cli2.JoinLobby(lobby)
	}, 6)

	assert.NoError(t, cli1.Close())
	assert.NoError(t, cli2.Close())

	fw.AssertSession(t)
}

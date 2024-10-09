package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/usecase"
	"github.com/stretchr/testify/assert"
)

var (
	ApiHost    = "http://localhost:8001/api/v1"
	SIOHost    = "http://localhost:8001"
	HttpClient = &http.Client{Timeout: 10 * time.Second}
	WaitTime   = 10 * time.Second
)

func PostUserWithID(t *testing.T, user *domain.User) *domain.User {
	b, err := json.Marshal(user)
	assert.NoError(t, err)

	resp, err := HttpClient.Post(fmt.Sprintf("%s/users/with_id", ApiHost), "application/json", bytes.NewReader(b))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	respUser := &domain.User{}
	err = json.NewDecoder(resp.Body).Decode(respUser)
	assert.NoError(t, err)
	return respUser
}

func FindLobby(t *testing.T) *domain.Lobby {
	findLobbyInput := usecase.FindLobbyInput{
		Dist: 0,
		// ИТМО - Кронверкский проспект, 49
		Location: domain.Coordinate{Lon: 30.310011, Lat: 59.956363},
	}
	b, err := json.Marshal(findLobbyInput)
	assert.NoError(t, err)

	resp, err := HttpClient.Post(fmt.Sprintf("%s/lobbies/find", ApiHost), "application/json", bytes.NewReader(b))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	lobby := &domain.Lobby{}
	assert.NoError(t, json.NewDecoder(resp.Body).Decode(lobby))
	return lobby
}

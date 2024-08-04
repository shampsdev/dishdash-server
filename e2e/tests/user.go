package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"dishdash.ru/internal/domain"

	"github.com/stretchr/testify/assert"
)

func UpdateUser(t *testing.T) {
	user := &domain.User{
		Name:     "name1",
		Avatar:   "avatar1",
		Telegram: nil,
	}
	b, err := json.Marshal(user)
	assert.NoError(t, err)

	// Post user
	user = postUser(t, user)

	resp, err := httpClient.Post(fmt.Sprintf("%s/users", ApiHost), "application/json", bytes.NewReader(b))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(t, err)

	// Get user
	resp, err = httpClient.Get(fmt.Sprintf("%s/users/%s", ApiHost, user.ID))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	respUser := &domain.User{}
	err = json.NewDecoder(resp.Body).Decode(respUser)
	assert.NoError(t, err)
	assertEqualUsers(t, user, respUser)

	// Update user
	user.Avatar = "new_avatar"
	b, err = json.Marshal(user)
	assert.NoError(t, err)
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/users", ApiHost),
		bytes.NewReader(b),
	)
	assert.NoError(t, err)
	resp, err = httpClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get user
	resp, err = httpClient.Get(fmt.Sprintf("%s/users/%s", ApiHost, user.ID))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(respUser)
	assert.NoError(t, err)
	assertEqualUsers(t, user, respUser)
}

func assertEqualUsers(t *testing.T, exp, actual *domain.User) {
	assert.Equal(t, exp.ID, actual.ID)
	assert.Equal(t, exp.Name, actual.Name)
	assert.Equal(t, exp.Avatar, actual.Avatar)
	assert.Equal(t, exp.Telegram, actual.Telegram)
}

func postUser(t *testing.T, user *domain.User) *domain.User {
	b, err := json.Marshal(user)
	assert.NoError(t, err)

	resp, err := httpClient.Post(fmt.Sprintf("%s/users", ApiHost), "application/json", bytes.NewReader(b))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	respUser := &domain.User{}
	err = json.NewDecoder(resp.Body).Decode(respUser)
	assert.NoError(t, err)
	return respUser
}

package tests

import (
	"bytes"
	"dishdash.ru/internal/domain"
	"encoding/json"
	"fmt"
	"gotest.tools/v3/assert"
	"net/http"
	"testing"
	"time"
)

func UpdateUser(t *testing.T, host string) {
	cli := http.Client{Timeout: 10 * time.Second}

	user := &domain.User{
		Name:     "name1",
		Avatar:   "avatar1",
		Telegram: nil,
	}
	b, err := json.Marshal(user)
	assert.NilError(t, err)

	// Post user
	resp, err := cli.Post(fmt.Sprintf("%s/users", host), "application/json", bytes.NewReader(b))
	assert.NilError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NilError(t, err)

	// Get user
	resp, err = cli.Get(fmt.Sprintf("%s/users/%s", host, user.ID))
	assert.NilError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	respUser := &domain.User{}
	err = json.NewDecoder(resp.Body).Decode(respUser)
	assert.NilError(t, err)
	assertEqualUsers(t, user, respUser)

	// Update user
	user.Avatar = "new_avatar"
	b, err = json.Marshal(user)
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/users", host),
		bytes.NewReader(b),
	)
	assert.NilError(t, err)
	resp, err = cli.Do(req)
	assert.NilError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get user
	resp, err = cli.Get(fmt.Sprintf("%s/users/%s", host, user.ID))
	assert.NilError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(respUser)
	assert.NilError(t, err)
	assertEqualUsers(t, user, respUser)
}

func assertEqualUsers(t *testing.T, actual *domain.User, exp *domain.User) {
	assert.Equal(t, actual.ID, exp.ID)
	assert.Equal(t, actual.Name, exp.Name)
	assert.Equal(t, actual.Avatar, exp.Avatar)
	assert.Equal(t, actual.Telegram, exp.Telegram)
}

package tests

import (
	"dishdash.ru/e2e/pg_test"
	"dishdash.ru/internal/domain"
	"encoding/json"
	"fmt"
	"gotest.tools/v3/assert"
	"net/http"
	"testing"
	"time"
)

func GetAllTags(t *testing.T, host string) {
	cli := http.Client{Timeout: 10 * time.Second}

	resp, err := cli.Get(fmt.Sprintf("%s/places/tags", host))
	assert.NilError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var tags []*domain.Tag
	err = json.NewDecoder(resp.Body).Decode(&tags)
	assert.NilError(t, err)
	assert.Equal(t, len(tags), len(pg_test.Tags))
	for i := range tags {
		assert.Equal(t, *tags[i], *pg_test.Tags[i])
	}
}

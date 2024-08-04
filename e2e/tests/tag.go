package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"dishdash.ru/e2e/pg_test"
	"dishdash.ru/internal/domain"

	"github.com/stretchr/testify/assert"
)

func GetAllTags(t *testing.T) {
	cli := http.Client{Timeout: 10 * time.Second}

	resp, err := cli.Get(fmt.Sprintf("%s/places/tags", ApiHost))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var tags []*domain.Tag
	err = json.NewDecoder(resp.Body).Decode(&tags)
	assert.NoError(t, err)
	assert.Equal(t, len(pg_test.Tags), len(tags))
	for i := range tags {
		assert.Equal(t, *pg_test.Tags[i], *tags[i])
	}
}

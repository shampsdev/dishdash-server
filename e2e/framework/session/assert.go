package session

import (
	"bytes"
	"encoding/json"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func AssertEqual(t *testing.T, exp, act *Session) {
	t.Helper()

	opts := cmp.Options{
		cmpopts.IgnoreUnexported(Session{}, Step{}),
		cmpopts.SortSlices(compareEventData),
		cmpopts.IgnoreMapEntries(isIgnoredMapEntry),
	}

	if !cmp.Equal(exp, act, opts) {
		t.Errorf("Sessions are not equal: \n%s", cmp.Diff(exp, act, opts))
	}
}

func compareEventData(x, y EventData) bool {
	aBytes, _ := json.Marshal(x)
	bBytes, _ := json.Marshal(y)
	return bytes.Compare(aBytes, bBytes) < 0
}

func isIgnoredMapEntry(key string, _ interface{}) bool {
	return slices.Contains([]string{
		"createdAt",
		"updatedAt",
	}, key)
}

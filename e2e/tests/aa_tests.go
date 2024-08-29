package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"slices"
	"sync"
	"testing"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	goassert "gotest.tools/v3/assert"
)

var (
	ApiHost    string
	SIOHost    string
	httpClient = &http.Client{Timeout: 10 * time.Second}
	waitTime   = 10 * time.Second
)

type SocketIOSession struct {
	lock       sync.Mutex
	UserEvents map[string]*eventCollection
}

type eventCollection struct {
	lock  sync.Mutex
	Steps []eventStep
}

type eventStep struct {
	Name   string
	Events []eventData
}

type eventData struct {
	Name string `json:"EventName"`
	Data map[string]interface{}
}

func newSocketIOSession() *SocketIOSession {
	return &SocketIOSession{
		UserEvents: make(map[string]*eventCollection),
	}
}

func (s *SocketIOSession) Save(filename string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s.UserEvents)
}

func LoadSocketIOSession(filename string) (*SocketIOSession, error) {
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	ec := &SocketIOSession{}
	if err := json.Unmarshal(byteValue, &ec.UserEvents); err != nil {
		return nil, err
	}
	return ec, nil
}

func (s *SocketIOSession) addUser(user string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.UserEvents[user] = &eventCollection{}
}

func (s *SocketIOSession) sioAddFunc(user, eventName string) func(socketio.Conn, map[string]interface{}) {
	return func(_ socketio.Conn, data map[string]interface{}) {
		s.add(user, eventData{
			Name: eventName,
			Data: data,
		})
	}
}

func (s *SocketIOSession) add(user string, event eventData) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.UserEvents[user].add(event)
}

func (s *SocketIOSession) newStep(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for k := range s.UserEvents {
		s.UserEvents[k].newStep(name)
	}
}

func AssertSocketIOSession(t *testing.T, exp, actual *SocketIOSession) {
	exp.lock.Lock()
	defer exp.lock.Unlock()
	actual.lock.Lock()
	defer actual.lock.Unlock()
	assert.Equal(t, len(exp.UserEvents), len(actual.UserEvents))
	for name := range exp.UserEvents {
		assertEventCollection(t, exp.UserEvents[name], actual.UserEvents[name])
	}
}

func (ec *eventCollection) newStep(name string) {
	ec.lock.Lock()
	defer ec.lock.Unlock()
	ec.Steps = append(ec.Steps, eventStep{Name: name, Events: []eventData{}})
}

func (ec *eventCollection) add(event eventData) {
	ec.lock.Lock()
	defer ec.lock.Unlock()
	ec.Steps[len(ec.Steps)-1].Events = append(ec.Steps[len(ec.Steps)-1].Events, event)
}

func assertEventCollection(t *testing.T, exp, actual *eventCollection) {
	t.Helper()
	exp.lock.Lock()
	defer exp.lock.Unlock()
	actual.lock.Lock()
	defer actual.lock.Unlock()

	t.Helper()
	for i := range exp.Steps {
		aStep := exp.Steps[i]
		bStep := actual.Steps[i]
		assertEventStep(t, aStep, bStep)
	}
}

func assertEventStep(t *testing.T, exp, actual eventStep) {
	t.Helper()
	assert.Equal(t, exp.Name, actual.Name)
	assert.Equal(t, len(exp.Events), len(actual.Events))

	eventDataCmp := func(a, b eventData) int {
		aBytes, err := json.Marshal(a)
		assert.NoError(t, err)
		bBytes, err := json.Marshal(b)
		assert.NoError(t, err)
		return bytes.Compare(aBytes, bBytes)
	}

	slices.SortFunc(exp.Events, eventDataCmp)
	slices.SortFunc(actual.Events, eventDataCmp)

	for i := range exp.Events {
		assert.Equal(t, exp.Events[i].Name, actual.Events[i].Name)
		assertMaps(t, exp.Events[i].Data, actual.Events[i].Data)
	}
}

var ignoredFields = []string{
	"updatedAt",
	"createdAt",
}

func assertMaps(t *testing.T, exp, actual map[string]interface{}) {
	var expCopy, actualCopy map[string]interface{}
	assert.NoError(t, copier.Copy(&expCopy, &exp))
	assert.NoError(t, copier.Copy(&actualCopy, &actual))
	removeDeepKeys(expCopy, ignoredFields)
	removeDeepKeys(actualCopy, ignoredFields)
	goassert.DeepEqual(t, expCopy, actualCopy)
}

func removeDeepKeys(m map[string]interface{}, keys []string) map[string]interface{} {
	for _, key := range keys {
		delete(m, key)
	}

	for key, value := range m {
		if reflect.ValueOf(value).Kind() == reflect.Map {
			if converted, ok := value.(map[string]interface{}); ok {
				m[key] = removeDeepKeys(converted, keys)
			}
		}
	}

	return m
}

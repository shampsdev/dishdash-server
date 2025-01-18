package state

import (
	"testing"

	"dishdash.ru/internal/domain"
	"github.com/stretchr/testify/assert"
)

type TState struct {
	id string

	lastSmthEvent *TEvent
	lastSmthCtx   *Context[*TState]
}

func (s *TState) ID() string {
	return s.id
}

type TEvent struct {
	S string `json:"s"`
	X int    `json:"x"`
}

func (s *TState) OnSomething1(c *Context[*TState]) {
	s.lastSmthCtx = c
	s.lastSmthEvent = nil
}

func (s *TState) OnSomething2(c *Context[*TState], ev TEvent) {
	s.lastSmthEvent = &ev
	s.lastSmthCtx = c
}

func TestContextFunctions(t *testing.T) {
	c := NewContext[*TState](nil, nil)
	c.State = &TState{id: "state1"}
	c.User = &domain.User{ID: "user1"}

	f1Calls := 0
	f1 := func(c *Context[*TState]) {
		assert.Equal(t, c.State.ID(), "state1")
		assert.Equal(t, c.User.ID, "user1")
		f1Calls++
	}

	assert.NoError(t, c.Call(f1))
	assert.Equal(t, f1Calls, 1)
	assert.NoError(t, c.Call(f1, struct{}{}))
	assert.Equal(t, f1Calls, 2)

	f2Calls := 0
	f2 := func(c *Context[*TState], ev TEvent) {
		assert.Equal(t, c.State.ID(), "state1")
		assert.Equal(t, c.User.ID, "user1")
		assert.Equal(t, ev.S, "hello")
		assert.Equal(t, ev.X, 1)
		f2Calls++
	}

	assert.Error(t, c.Call(f2))
	assert.Equal(t, f2Calls, 0)

	assert.NoError(t, c.Call(f2, TEvent{S: "hello", X: 1}))
	assert.Equal(t, f2Calls, 1)

	assert.NoError(t, c.Call(f2, map[string]interface{}{"s": "hello", "x": 1}))
	assert.Equal(t, f2Calls, 2)
}

func TestContextStruct(t *testing.T) {
	c := NewContext[*TState](nil, nil)
	c.State = &TState{id: "state1"}
	c.User = &domain.User{ID: "user1"}

	assert.NoError(t, c.Call((*TState).OnSomething2, TEvent{S: "hello", X: 1}))
	assert.Equal(t, c.State.lastSmthEvent.S, "hello")
	assert.Equal(t, c.State.lastSmthEvent.X, 1)
	assert.Equal(t, c.State.lastSmthCtx.State.ID(), "state1")
	assert.Equal(t, c.State.lastSmthCtx.User.ID, "user1")

	assert.NoError(t, c.Call((*TState).OnSomething1))
	assert.Nil(t, c.State.lastSmthEvent)
	assert.Equal(t, c.State.lastSmthCtx.State.ID(), "state1")
	assert.Equal(t, c.State.lastSmthCtx.User.ID, "user1")
}

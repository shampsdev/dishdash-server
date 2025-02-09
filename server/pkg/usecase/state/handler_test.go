package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TState struct {
	id            string
	lastSmthEvent *TEvent
	lastSmthCtx   *Context[*TState]
}

func (s *TState) ID() string {
	return s.id
}

type TEvent struct {
	Hello string `json:"hello"`
	X     int    `json:"x"`
}

func (s *TState) OnSomething(c *Context[*TState], ev TEvent) error {
	s.lastSmthEvent = &ev
	s.lastSmthCtx = c
	return nil
}

func TestWrapMethod(t *testing.T) {
	s := &TState{}
	c := NewContext[*TState](nil, nil)
	c.State = s

	wrappedMethod := WrapHMethod((*TState).OnSomething)
	assert.NoError(t, wrappedMethod(c, TEvent{Hello: "world", X: 42}))
	assert.Equal(t, s.lastSmthEvent, &TEvent{Hello: "world", X: 42})
	assert.Equal(t, s.lastSmthCtx, c)

	assert.NoError(t, wrappedMethod(c, map[string]interface{}{"hello": "world", "x": 42}))
	assert.Equal(t, s.lastSmthEvent, &TEvent{Hello: "world", X: 42})
	assert.Equal(t, s.lastSmthCtx, c)
}

func TestWrapTyped(t *testing.T) {
	s := &TState{}
	c := NewContext[*TState](nil, nil)
	c.State = s

	f := func(c *Context[*TState], ev TEvent) error {
		s.lastSmthEvent = &ev
		s.lastSmthCtx = c
		return nil
	}

	wrappedMethod := WrapHTyped(f)
	assert.NoError(t, wrappedMethod(c, TEvent{Hello: "world", X: 42}))
	assert.Equal(t, s.lastSmthEvent, &TEvent{Hello: "world", X: 42})
	assert.Equal(t, s.lastSmthCtx, c)

	assert.NoError(t, wrappedMethod(c, nil))
	assert.Equal(t, s.lastSmthEvent, &TEvent{})
	assert.Equal(t, s.lastSmthCtx, c)
}

package state

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type HandlerFunc[State WithID] func(c *Context[State], arg interface{}) error

type HandlerFuncTyped[State WithID, Arg any] func(c *Context[State], arg Arg) error

type HandlerFuncMethod[State WithID, Arg any] func(s State, c *Context[State], arg Arg) error

func WrapHTyped[State WithID, Arg any](
	f HandlerFuncTyped[State, Arg],
) HandlerFunc[State] {
	return func(c *Context[State], arg interface{}) error {
		var argTyped Arg
		if arg == nil {
			return f(c, argTyped)
		}

		err := mapstructure.Decode(arg, &argTyped)
		if err != nil {
			return fmt.Errorf("failed to decode event data: %w", err)
		}

		return f(c, argTyped)
	}
}

func WrapHMethod[State WithID, Arg any](
	f HandlerFuncMethod[State, Arg],
) HandlerFunc[State] {
	return WrapHTyped(
		func(c *Context[State], arg Arg) error {
			return f(c.State, c, arg)
		},
	)
}

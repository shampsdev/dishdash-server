package state

type WithID interface {
	ID() string
}

type Server[State WithID] interface {
	ForEach(stateID string, f func(c *Context[State]))
}

type Conn interface {
	Emit(event string, data interface{})
	Close() error
}

package state

import (
	"fmt"
	rf "reflect"

	"github.com/mitchellh/mapstructure"
)

// Call can call these types of function:
// func(Context)
// func(Context, any)
// func(State, Context)
// func(State, Context, any)
func (c *Context[State]) Call(f interface{}, args ...interface{}) error {
	var eventDataCasted any

	if len(args) == 1 && rf.TypeOf(f).In(rf.TypeOf(f).NumIn()-1).Kind() != rf.TypeOf(c).Kind() {
		eventDataCasted = rf.New(rf.TypeOf(f).In(rf.TypeOf(f).NumIn() - 1)).Elem().Interface()
		err := mapstructure.Decode(args[0], &eventDataCasted)
		if err != nil {
			return err
		}
	}

	if len(args) > 1 {
		return fmt.Errorf("too many arguments")
	}

	fNumIn := rf.ValueOf(f).Type().NumIn()
	fIn0Kind := rf.ValueOf(f).Type().In(0).Kind()

	var vArgs []rf.Value

	if fNumIn == 1 {
		// func(Context)
		vArgs = []rf.Value{rf.ValueOf(c)}
	} else if fNumIn == 2 && fIn0Kind == rf.TypeOf(c).Kind() {
		// func(Context, any)
		vArgs = []rf.Value{rf.ValueOf(c), rf.ValueOf(eventDataCasted)}
	} else if fNumIn == 2 &&
		fIn0Kind == rf.TypeOf(c.State).Kind() {
		// func(State, Context)
		vArgs = []rf.Value{rf.ValueOf(c.State), rf.ValueOf(c)}
	} else if fNumIn == 3 &&
		fIn0Kind == rf.TypeOf(c.State).Kind() {
		// func(State, Context, any)
		vArgs = []rf.Value{rf.ValueOf(c.State), rf.ValueOf(c), rf.ValueOf(eventDataCasted)}
	} else {
		panic(fmt.Sprintf("illegal event handler %v %v", rf.TypeOf(f), rf.ValueOf(f)))
	}

	rets := rf.ValueOf(f).Call(vArgs)
	if len(rets) == 1 {
		err, ok := rets[0].Interface().(error)
		if ok && err != nil {
			return err
		}
	}
	return nil
}

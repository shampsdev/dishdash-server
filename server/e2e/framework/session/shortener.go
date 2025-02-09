package session

import (
	"encoding/json"
)

type Shortener func(any) any

func EventShortener[T any](f func(T) any) Shortener {
	return func(a any) any {
		var t T
		bs, err := json.Marshal(a)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(bs, &t); err != nil {
			panic(err)
		}
		return f(t)
	}
}

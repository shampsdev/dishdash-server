package room

import "context"

type Repo interface {
	GetRoom(ctx context.Context, id string) (*Room, error)
	DeleteRoom(ctx context.Context, id string) error
}

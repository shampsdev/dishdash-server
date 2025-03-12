package usecase

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Event struct {
	db *pgxpool.Pool
}

func NewEventUsecase(db *pgxpool.Pool) *Event {
	return &Event{
		db: db,
	}
}

type SaveEventInput struct {
	Name string
	Data json.RawMessage
}

func (u *Event) SaveEvent(ctx context.Context, event SaveEventInput) error {
	_, err := u.db.Exec(
		context.Background(),
		"INSERT INTO event (name, data) VALUES ($1, $2)",
		event.Name,
		event.Data,
	)
	return err
}

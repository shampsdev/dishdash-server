package usecase

import "github.com/jackc/pgx/v5/pgxpool"

type Cases struct {
	Event *Event
}

func Setup(db *pgxpool.Pool) Cases {
	return Cases{
		Event: NewEventUsecase(db),
	}
}

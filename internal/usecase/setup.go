package usecase

import (
	"dishdash.ru/internal/repo/pg"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(pool *pgxpool.Pool) Cases {
	pr := pg.NewPlaceRepo(pool)
	tr := pg.NewTagRepo(pool)
	lr := pg.NewLobbyRepo(pool)
	ur := pg.NewUserRepo(pool)
	sr := pg.NewSwipeRepo(pool)

	pu := NewPlaceUseCase(tr, pr)
	lu := NewLobbyUseCase(lr, ur, tr, pr, sr)

	return Cases{
		Place:    pu,
		Tag:      NewTagUseCase(tr),
		Lobby:    lu,
		User:     NewUserUseCase(ur),
		Swipe:    NewSwipeUseCase(sr),
		RoomRepo: NewInMemoryRoomRepo(lu, pu),
	}
}

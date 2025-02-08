package usecase

import (
	"dishdash.ru/pkg/repo/pg"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup(pool *pgxpool.Pool) Cases {
	pr := pg.NewPlaceRepo(pool)
	tr := pg.NewTagRepo(pool)
	lr := pg.NewLobbyRepo(pool)
	ur := pg.NewUserRepo(pool)
	sr := pg.NewSwipeRepo(pool)
	prr := pg.NewPlaceRecommenderRepo(pool)

	pu := NewPlaceUseCase(tr, pr)
	lu := NewLobbyUseCase(lr, ur, tr, pr, sr)
	su := NewSwipeUseCase(sr)
	uu := NewUserUseCase(ur)

	placeRecommender := NewPlaceRecommender(
		prr,
		pr,
		tr,
	)

	return Cases{
		Place:    pu,
		Tag:      NewTagUseCase(tr),
		Lobby:    lu,
		User:     uu,
		Swipe:    su,
		RoomRepo: NewInMemoryRoomRepo(lu, pu, su, uu, placeRecommender),
	}
}

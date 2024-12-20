package usecase

import (
	"dishdash.ru/cmd/server/config"
	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo/pg"
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
		domain.RecommendOpts{
			PriceCoeff: float64(config.C.Recommendation.PriceCoeff),
			DistCoeff:  float64(config.C.Recommendation.DistCoeff),
		},
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

package pg

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"

	"dishdash.ru/internal/domain"
)

const saveSwipeQuery = `
    INSERT INTO "swipe" (
        "lobby_id",
        "card_id",
        "user_id",
        "swipe_type"
    ) VALUES ($1, $2, $3, $4)
`

type SwipeRepository struct {
	db *pgxpool.Pool
}

func NewSwipeRepository(pool *pgxpool.Pool) *SwipeRepository {
	return &SwipeRepository{db: pool}
}

func (sr *SwipeRepository) SaveSwipe(ctx context.Context, swipe *domain.Swipe) error {
	_, err := sr.db.Exec(ctx, saveSwipeQuery,
		swipe.LobbyID,
		swipe.CardID,
		swipe.UserID,
		swipe.SwipeType,
	)
	if err != nil {
		log.Printf("Error saving swipe: %v\n", err)
		return err
	}
	return nil
}

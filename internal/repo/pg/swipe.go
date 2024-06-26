package pg

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"dishdash.ru/internal/domain"
)

type SwipeRepository struct {
	db *pgxpool.Pool
}

func NewSwipeRepository(pool *pgxpool.Pool) *SwipeRepository {
	return &SwipeRepository{db: pool}
}

func (sr *SwipeRepository) CreateSwipe(ctx context.Context, swipe *domain.Swipe) error {
	const saveSwipeQuery = `
	    INSERT INTO "swipe" (
	        "lobby_id",
	        "card_id",
	        "user_id",
	        "type"
	    ) VALUES ($1, $2, $3, $4)
    `

	_, err := sr.db.Exec(ctx, saveSwipeQuery,
		swipe.LobbyID,
		swipe.CardID,
		swipe.UserID,
		swipe.Type,
	)

	if err != nil {
		log.Printf("Error saving swipe: %v\n", err)
		return err
	}
	return nil
}

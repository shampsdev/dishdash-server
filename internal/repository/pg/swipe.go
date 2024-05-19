package pg

import (
	"context"
	"log"

	"dishdash.ru/internal/domain"
	"github.com/jackc/pgx/v4"
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
	db *pgx.Conn
}

func NewSwipeRepository(db *pgx.Conn) *SwipeRepository {
	return &SwipeRepository{db: db}
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

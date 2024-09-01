package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"dishdash.ru/internal/domain"
)

type SwipeRepo struct {
	db *pgxpool.Pool
}

func NewSwipeRepo(pool *pgxpool.Pool) *SwipeRepo {
	return &SwipeRepo{db: pool}
}

func (sr *SwipeRepo) SaveSwipe(ctx context.Context, swipe *domain.Swipe) error {
	const saveSwipeQuery = `
	    INSERT INTO "swipe" (lobby_id, place_id, user_id, type) 
	    VALUES ($1, $2, $3, $4)
    `

	_, err := sr.db.Exec(ctx, saveSwipeQuery,
		swipe.LobbyID,
		swipe.PlaceID,
		swipe.UserID,
		swipe.Type,
	)
	if err != nil {
		return fmt.Errorf("can't save swipe: %w", err)
	}
	return nil
}

func (sr *SwipeRepo) GetSwipesByLobbyID(ctx context.Context, lobbyID string) ([]*domain.Swipe, error) {
	const query = `
	SELECT lobby_id, place_id, user_id, type
	FROM "swipe" WHERE lobby_id = $1
`
	rows, err := sr.db.Query(ctx, query, lobbyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	swipes := make([]*domain.Swipe, 0)
	for rows.Next() {
		var swipe domain.Swipe
		err := rows.Scan(
			&swipe.LobbyID,
			&swipe.PlaceID,
			&swipe.UserID,
			&swipe.Type,
		)
		if err != nil {
			return nil, err
		}
	}
	return swipes, nil
}

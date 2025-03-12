package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"dishdash.ru/pkg/domain"
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
		swipe.CardID,
		swipe.UserID,
		swipe.Type,
	)
	if err != nil {
		return fmt.Errorf("can't save swipe: %w", err)
	}
	return nil
}

func (sr *SwipeRepo) GetSwipesCount(ctx context.Context) (int, error) {
	const query = `
	SELECT COUNT(*)
	FROM "swipe"
`
	row := sr.db.QueryRow(ctx, query)
	var count int
	err := row.Scan(
		&count)
	if err != nil {
		return 0, err
	}

	return count, err
}

func (sr *SwipeRepo) GetSwipesByLobbyID(ctx context.Context, lobbyID string) ([]*domain.Swipe, error) {
	const query = `
	SELECT lobby_id, place_id, user_id, type
	FROM "swipe" WHERE lobby_id = $1 ORDER BY id
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
			&swipe.CardID,
			&swipe.UserID,
			&swipe.Type,
		)
		if err != nil {
			return nil, err
		}
		swipes = append(swipes, &swipe)
	}
	return swipes, nil
}

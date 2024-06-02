package pg

import (
	"context"
	"dishdash.ru/internal/domain"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TagRepository struct {
	db *pgxpool.Pool
}

func NewTagRepository(db *pgxpool.Pool) *TagRepository {
	return &TagRepository{db: db}
}

func (t *TagRepository) SaveTag(ctx context.Context, tag *domain.Tag) error {
	query := `INSERT INTO tag (name, icon) VALUES ($1, $2) RETURNING id`
	err := t.db.QueryRow(ctx, query, tag.Name, tag.Icon).Scan(&tag.ID)
	if err != nil {
		return fmt.Errorf("could not insert tag: %w", err)
	}
	return nil
}

func (t *TagRepository) AttachTagToCard(ctx context.Context, tagID int64, cardID int64) error {
	query := `INSERT INTO tag_card (tag_id, card_id) VALUES ($1, $2)`
	_, err := t.db.Exec(ctx, query, tagID, cardID)
	if err != nil {
		return fmt.Errorf("could not attach tag to card: %w", err)
	}
	return nil
}

func (t *TagRepository) GetTagsByCardID(ctx context.Context, cardID int64) ([]*domain.Tag, error) {
	query := `
	SELECT t.id, t.name, t.icon
	FROM tag t
	JOIN tag_card tc ON t.id = tc.tag_id
	WHERE tc.card_id = $1
	`

	rows, err := t.db.Query(ctx, query, cardID)
	if err != nil {
		return nil, fmt.Errorf("could not get tags by card ID: %w", err)
	}
	defer rows.Close()

	var tags []*domain.Tag
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Icon)
		if err != nil {
			return nil, fmt.Errorf("could not scan tag: %w", err)
		}
		tags = append(tags, &tag)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return tags, nil
}

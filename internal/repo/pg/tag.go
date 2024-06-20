package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"dishdash.ru/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TagRepository struct {
	db *pgxpool.Pool
}

func NewTagRepository(db *pgxpool.Pool) *TagRepository {
	return &TagRepository{db: db}
}

func (t *TagRepository) CreateTag(ctx context.Context, tag *domain.Tag) (int64, error) {
	query := `INSERT INTO tag (name, icon) VALUES ($1, $2) RETURNING id`
	var id int64
	err := t.db.QueryRow(ctx, query, tag.Name, tag.Icon).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("could not insert tag: %w", err)
	}
	return id, err
}

func (t *TagRepository) AttachTagsToCard(ctx context.Context, tagIDs []int64, cardID int64) error {
	batch := &pgx.Batch{}

	query := `INSERT INTO card_tag (tag_id, card_id) VALUES ($1, $2)`
	for _, tagID := range tagIDs {
		batch.Queue(query, tagID, cardID)
	}

	br := t.db.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		return fmt.Errorf("could not attach tags to card: %w", err)
	}
	return nil
}

func (t *TagRepository) GetTagsByCardID(ctx context.Context, cardID int64) ([]*domain.Tag, error) {
	query := `
	SELECT tag.id, tag.name, tag.icon
	FROM tag
	JOIN card_tag ON tag.id = card_tag.tag_id
	WHERE card_tag.card_id = $1
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

func (t *TagRepository) GetAllTags(ctx context.Context) ([]*domain.Tag, error) {
	query := `
	SELECT tag.id, tag.name, tag.icon
	FROM tag
	`

	rows, err := t.db.Query(ctx, query)
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

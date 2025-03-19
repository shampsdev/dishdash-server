package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"dishdash.ru/pkg/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TagRepo struct {
	db *pgxpool.Pool
}

func NewTagRepo(db *pgxpool.Pool) *TagRepo {
	return &TagRepo{db: db}
}

func (tr *TagRepo) SaveTag(ctx context.Context, tag *domain.Tag) (int64, error) {
	query := `INSERT INTO "tag" ("name", "icon", "visible", "order", "excluded") VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	err := tr.db.QueryRow(ctx, query, tag.Name, tag.Icon, tag.Visible, tag.Order, tag.Excluded).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("could not insert tag: %w", err)
	}
	return id, err
}

func (tr *TagRepo) DeleteTag(ctx context.Context, tagId int64) error {
	query := `DELETE FROM "tag" WHERE id=$1`
	_, err := tr.db.Exec(ctx, query, tagId)
	if err != nil {
		return fmt.Errorf("could not delete tag: %w", err)
	}
	return nil
}

func (tr *TagRepo) UpdateTag(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	query := `UPDATE "tag" SET "name" = $1, "icon" = $2, "visible" = $3, "order" = $4, "excluded" = $5 WHERE id = $5`
	_, err := tr.db.Exec(ctx, query, tag.Name, tag.Icon, tag.Visible, tag.Order, tag.Excluded, tag.ID)
	if err != nil {
		return tag, fmt.Errorf("could not update tag: %w", err)
	}
	return tag, nil
}

func (tr *TagRepo) AttachTagsToPlace(ctx context.Context, tagIDs []int64, placeID int64) error {
	if len(tagIDs) == 0 {
		return nil
	}
	batch := &pgx.Batch{}

	query := `INSERT INTO "place_tag" ("tag_id", "place_id") 
        VALUES ($1, $2)
        ON CONFLICT ("place_id", "tag_id") DO NOTHING`
	for _, tagID := range tagIDs {
		batch.Queue(query, tagID, placeID)
	}

	br := tr.db.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		return fmt.Errorf("could not attach tags to place: %w", err)
	}
	return nil
}

func (tr *TagRepo) DetachTagsFromPlace(ctx context.Context, placeID int64) error {
	query := `DELETE FROM "place_tag" WHERE "place_id" = $1`
	_, err := tr.db.Exec(ctx, query, placeID)
	if err != nil {
		return fmt.Errorf("could not detach tags from place: %w", err)
	}
	return nil
}

func (tr *TagRepo) DetachTagsFromLobby(ctx context.Context, lobbyID string) error {
	query := `DELETE FROM "lobby_tag" WHERE "lobby_id" = $1`
	_, err := tr.db.Exec(ctx, query, lobbyID)
	if err != nil {
		return fmt.Errorf("could not detach tags from lobby: %w", err)
	}
	return nil
}

func (tr *TagRepo) AttachTagsToLobby(ctx context.Context, tagIDs []int64, lobbyID string) error {
	if len(tagIDs) == 0 {
		return nil
	}
	batch := &pgx.Batch{}

	query := `INSERT INTO "lobby_tag" ("tag_id", "lobby_id") VALUES ($1, $2)`
	for _, tagID := range tagIDs {
		batch.Queue(query, tagID, lobbyID)
	}

	br := tr.db.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		return fmt.Errorf("could not attach tags to lobby: %w", err)
	}
	return nil
}

func (tr *TagRepo) GetTagsByPlaceID(ctx context.Context, placeID int64) ([]*domain.Tag, error) {
	query := `
	SELECT tag.id, tag.name, tag.icon, tag.visible, tag.order, tag.excluded
	FROM tag
	JOIN place_tag ON tag.id = place_tag.tag_id
	WHERE place_tag.place_id = $1
	`

	rows, err := tr.db.Query(ctx, query, placeID)
	if err != nil {
		return nil, fmt.Errorf("could not get tags by place ID: %w", err)
	}
	defer rows.Close()

	tags := make([]*domain.Tag, 0)
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Icon, &tag.Visible, &tag.Order, &tag.Excluded)
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

func (tr *TagRepo) GetTagsByLobbyID(ctx context.Context, lobbyID string) ([]*domain.Tag, error) {
	query := `
	SELECT tag.id, tag.name, tag.icon, tag.visible, tag.order, tag.excluded
	FROM tag
	JOIN lobby_tag ON tag.id = lobby_tag.tag_id
	WHERE lobby_tag.lobby_id = $1
	`

	rows, err := tr.db.Query(ctx, query, lobbyID)
	if err != nil {
		return nil, fmt.Errorf("could not get tags by lobby ID: %w", err)
	}
	defer rows.Close()

	tags := make([]*domain.Tag, 0)
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Icon, &tag.Visible, &tag.Order, &tag.Excluded)
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

func (tr *TagRepo) GetAllTags(ctx context.Context) ([]*domain.Tag, error) {
	query := `
	SELECT tag.id, tag.name, tag.icon, tag.visible, tag.order, tag.excluded
	FROM tag
	`

	rows, err := tr.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get tags by place ID: %w", err)
	}
	defer rows.Close()

	tags := make([]*domain.Tag, 0)
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Icon, &tag.Visible, &tag.Order, &tag.Excluded)
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

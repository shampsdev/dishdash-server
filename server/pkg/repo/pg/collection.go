package pg

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"

	"dishdash.ru/pkg/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CollectionRepo struct {
	db *pgxpool.Pool
}

func (cr *CollectionRepo) SaveCollection(ctx context.Context, collection *domain.Collection) (int64, error) {
	const insertCollectionQuery = `
		INSERT INTO "collection" (name, description)
		VALUES ($1, $2)
		RETURNING id;
	`

	tx, err := cr.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("rollback failed: %v", err)
		}
	}()

	var collectionID int64
	err = tx.QueryRow(ctx, insertCollectionQuery, collection.Name, collection.Description).Scan(&collectionID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert collection: %w", err)
	}

	if len(collection.Places) > 0 {
		const insertCollectionPlaceQuery = `
			INSERT INTO "collection_place" (collection_id, place_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING;
		`

		for _, place := range collection.Places {
			_, err := tx.Exec(ctx, insertCollectionPlaceQuery, collectionID, place.ID)
			if err != nil {
				return 0, fmt.Errorf("failed to link place %d with collection %d: %w", place.ID, collectionID, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return collectionID, nil
}

func (cr *CollectionRepo) GetCollectionByID(ctx context.Context, collectionID int64) (*domain.Collection, error) {
	const getCollectionQuery = `
	SELECT
    c.id,
    c.name,
    c.description,
	c.created_at,
	c.updated_at,
	c.avatar,
	c.visible,
	c.order,
    COALESCE(
        JSON_AGG(
            JSON_BUILD_OBJECT(
                'id', p.id,
                'name', p.title,
                'short_description', p.short_description,
                'description', p.description,
                'images', p.images,
                'location', p.location,
                'address', p.address,
                'priceAvg', p.price_avg,
                'reviewRating', p.review_ratin,
                'reviewCount', p.review_count,
                'updatedAt', p.updated_at,
                'source', p.source,
                'url', p.url,
                'boost', p.boost,
                'boostRadius', p.boost_radius,
            )
        ) FILTER (WHERE p.id IS NOT NULL)
    ) AS places
FROM "collection" AS c
LEFT JOIN "collection_place" AS cp ON c.id = cp.collection_id 
LEFT JOIN "place" AS p ON cp.place_id = p.id
WHERE c.id = $1
GROUP BY c.id;

`

	row := cr.db.QueryRow(ctx, getCollectionQuery, collectionID)
	collection, err := scanCollection(row)
	if err != nil {
		return nil, fmt.Errorf("can't fetch collection: %w", err)
	}
	return collection, nil
}

func (cr *CollectionRepo) GetAllCollections(ctx context.Context) ([]*domain.Collection, error) {
	const getCollectionQuery = `
	SELECT
	    c.id,
	    c.name,
	    c.description,
	    c.created_at,
		c.updated_at,
		c.avatar,
		c.visible,
		c.order,
	    COALESCE(
	        JSON_AGG(
	            JSON_BUILD_OBJECT(
	                'id', p.id,
	                'name', p.title,
	                'short_description', p.short_description,
	                'description', p.description,
	                'images', p.images,
	                'location', p.location,
	                'address', p.address,
	                'priceAvg', p.price_avg,
	                'reviewRating', p.review_rating,
	                'reviewCount', p.review_count,
	                'updatedAt', p.updated_at,
	                'source', p.source,
	                'url', p.url,
	                'boost', p.boost,
	                'boostRadius', p.boost_radius
	            )
	        ) FILTER (WHERE p.id IS NOT NULL)
	    ) AS places
	FROM "collection" AS c
	LEFT JOIN "collection_place" AS cp ON c.id = cp.collection_id 
	LEFT JOIN "place" AS p ON cp.place_id = p.id
	GROUP BY c.id;
	`

	rows, err := cr.db.Query(ctx, getCollectionQuery)
	if err != nil {
		return nil, fmt.Errorf("can't fetch collections: %w", err)
	}
	defer rows.Close()

	var collections []*domain.Collection

	for rows.Next() {
		collection, err := scanCollection(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning collection: %w", err)
		}
		collections = append(collections, collection)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over collections: %w", err)
	}

	return collections, nil
}

func (cr *CollectionRepo) DeleteCollectionByID(ctx context.Context, collectionID int64) error {
	const deleteCollectionQuery = `DELETE FROM "collection" WHERE id = $1;`

	_, err := cr.db.Exec(ctx, deleteCollectionQuery, collectionID)
	if err != nil {
		return fmt.Errorf("failed to delete collection %d: %w", collectionID, err)
	}

	return nil
}

func (cr *CollectionRepo) AttachPlacesToCollection(ctx context.Context, placeIDs []int64, collectionID int64) error {
	if len(placeIDs) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	const query = `
		INSERT INTO "collection_places" (collection_id, place_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
`
	for _, placeID := range placeIDs {
		batch.Queue(query, collectionID, placeID)
	}

	br := cr.db.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		return fmt.Errorf("could not attach places to collection: %w", err)
	}
	return nil
}

func (cr *CollectionRepo) DetachPlaceFromCollection(ctx context.Context, placeID, collectionID int64) error {
	const deleteQuery = `
		DELETE FROM "collection_places"
		WHERE collection_id = $1 AND place_id = $2;
	`

	_, err := cr.db.Exec(ctx, deleteQuery, collectionID, placeID)
	if err != nil {
		return fmt.Errorf("failed to detach place %d from collection %d: %w", placeID, collectionID, err)
	}

	return nil
}

func scanCollection(s Scanner) (*domain.Collection, error) {
	collection := new(domain.Collection)
	placesJSON := ""

	err := s.Scan(
		&collection.ID,
		&collection.Name,
		&collection.Description,
		&collection.CreatedAt,
		&collection.UpdatedAt,
		&collection.Avatar,
		&collection.Visible,
		&collection.Order,
		&placesJSON,
	)
	if err != nil {
		return nil, err
	}

	var places []domain.Place
	if err := json.Unmarshal([]byte(placesJSON), &places); err != nil {
		return nil, fmt.Errorf("failed to unmarshal places: %w", err)
	}

	collection.Places = make([]*domain.Place, len(places))
	for i := range places {
		collection.Places[i] = &places[i]
	}

	return collection, nil
}

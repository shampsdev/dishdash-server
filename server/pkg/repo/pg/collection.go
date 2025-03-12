package pg

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"dishdash.ru/pkg/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CollectionRepo struct {
	db   *pgxpool.Pool
	rand *rand.Rand
}

func NewCollectionRepo(db *pgxpool.Pool) *CollectionRepo {
	return &CollectionRepo{db: db, rand: rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), rand.Uint64()))}
}

func (cr *CollectionRepo) SaveCollection(ctx context.Context, collection *domain.Collection) (string, error) {
	const saveQuery = `
	INSERT INTO "collection" (
		"id",
        "name",
        "description",
        "avatar",
        "visible",
        "updated_at", "created_at")
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING "id"
		`

	collection.ID = cr.generateID()
	collection.UpdatedAt = time.Now().UTC()
	collection.CreatedAt = time.Now().UTC()
	row := cr.db.QueryRow(ctx, saveQuery,
		collection.ID,
		collection.Name,
		collection.Description,
		collection.Avatar,
		collection.Visible,
		collection.UpdatedAt,
		collection.CreatedAt,
	)
	var id string

	err := row.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (cr *CollectionRepo) GetCollectionByID(ctx context.Context, collectionID string) (*domain.Collection, error) {
	const getCollectionQuery = `
	SELECT 
		id,
		name,
		description,
		avatar,
		visible,
		updated_at,
		created_at,
		"order"
	FROM "collection"
	WHERE id = $1;
`

	row := cr.db.QueryRow(ctx, getCollectionQuery, collectionID)
	collection := new(domain.Collection)
	err := row.Scan(
		&collection.ID,
		&collection.Name,
		&collection.Description,
		&collection.Avatar,
		&collection.Visible,
		&collection.UpdatedAt,
		&collection.CreatedAt,
		&collection.Order,
	)
	if err != nil {
		return nil, fmt.Errorf("can't fetch collection: %w", err)
	}
	return collection, nil
}

func (cr *CollectionRepo) GetAllCollections(ctx context.Context) ([]*domain.Collection, error) {
	const getCollectionsQuery = `
        SELECT 
            id,
            name,
            description,
            avatar,
            visible,
            updated_at,
            created_at,
            "order"
        FROM "collection";
    `

	rows, err := cr.db.Query(ctx, getCollectionsQuery)
	if err != nil {
		return nil, fmt.Errorf("can't fetch collections: %w", err)
	}
	defer rows.Close()

	collections := make([]*domain.Collection, 0)
	for rows.Next() {
		collection := new(domain.Collection)
		err := rows.Scan(
			&collection.ID,
			&collection.Name,
			&collection.Description,
			&collection.Avatar,
			&collection.Visible,
			&collection.UpdatedAt,
			&collection.CreatedAt,
			&collection.Order,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan collection: %w", err)
		}
		collections = append(collections, collection)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning collections: %w", err)
	}

	return collections, nil
}

func (cr *CollectionRepo) GetPlacesByCollectionID(ctx context.Context, collectionID string) ([]*domain.Place, error) {
	const getPlacesQuery = `
		SELECT
			p.id,
			p.title,
			p.short_description,
			p.description,
			p.images,
			p.location,
			p.address,
			p.price_avg,
			p.review_rating,
			p.review_count,
			p.updated_at, 
			p.source,
			p.url,
			p.boost,
			p.boost_radius,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT('id', t.id, 'name', t.name, 'icon', t.icon, 'visible', t.visible, 'order', t.order)
				) FILTER (WHERE t.id IS NOT NULL),
				'[]'
			) AS tags
		FROM "place" AS p
		LEFT JOIN "place_tag" AS pt ON p.id = pt.place_id
		LEFT JOIN "tag" AS t ON pt.tag_id = t.id
		JOIN "collection_place" AS cp ON p.id = cp.place_id
		WHERE cp.collection_id = $1
		GROUP BY p.id;
    `

	rows, err := cr.db.Query(ctx, getPlacesQuery, collectionID)
	if err != nil {
		return nil, fmt.Errorf("can't fetch places: %w", err)
	}
	defer rows.Close()

	places := make([]*domain.Place, 0)
	for rows.Next() {
		place, err := scanPlace(rows)
		if err != nil {
			return nil, fmt.Errorf("could not get places by lobby ID: %w", err)
		}
		places = append(places, place)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning places: %w", err)
	}

	return places, nil
}

func (cr *CollectionRepo) GetAllCollectionsWithPlaces(ctx context.Context) ([]*domain.Collection, error) {
	collections, err := cr.GetAllCollections(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get collections: %w", err)
	}

	for _, collection := range collections {
		places, err := cr.GetPlacesByCollectionID(ctx, collection.ID)
		if err != nil {
			return nil, fmt.Errorf("can't get places for collection %s: %w", collection.ID, err)
		}
		collection.Places = places
	}

	return collections, nil
}

func (cr *CollectionRepo) GetCollectionWithPlacesByID(ctx context.Context, collectionID string) (*domain.Collection, error) {
	collection, err := cr.GetCollectionByID(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("can't get collection: %w", err)
	}

	places, err := cr.GetPlacesByCollectionID(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("can't get places for collection %s: %w", collectionID, err)
	}
	collection.Places = places

	return collection, nil
}

func (cr *CollectionRepo) DeleteCollectionByID(ctx context.Context, collectionID string) error {
	const deleteQuery = `
        DELETE FROM "collection" WHERE id = $1;
    `
	_, err := cr.db.Exec(ctx, deleteQuery, collectionID)
	if err != nil {
		return fmt.Errorf("can't delete collection: %w", err)
	}
	return nil
}

func (cr *CollectionRepo) AttachPlacesToCollection(ctx context.Context, placeIDs []int64, collectionID string) error {
	const attachQuery = `
        INSERT INTO "collection_place" (collection_id, place_id)
        VALUES ($1, $2);
    `

	for _, placeID := range placeIDs {
		_, err := cr.db.Exec(ctx, attachQuery, collectionID, placeID)
		if err != nil {
			return fmt.Errorf("can't attach place to collection: %w", err)
		}
	}
	return nil
}

func (cr *CollectionRepo) DetachPlacesFromCollection(ctx context.Context, collectionID string) error {
	const detachQuery = `
        DELETE FROM "collection_place" WHERE collection_id = $1;
    `
	_, err := cr.db.Exec(ctx, detachQuery, collectionID)
	if err != nil {
		return fmt.Errorf("can't detach places from collection: %w", err)
	}
	return nil
}

func (cr *CollectionRepo) UpdateCollection(ctx context.Context, collection *domain.Collection) error {
	const updateQuery = `
        UPDATE "collection"
        SET name=$1, description=$2, avatar=$3, visible=$4, updated_at=$5
        WHERE id=$6;
    `

	collection.UpdatedAt = time.Now().UTC()
	_, err := cr.db.Exec(ctx, updateQuery,
		collection.Name,
		collection.Description,
		collection.Avatar,
		collection.Visible,
		collection.UpdatedAt,
		collection.ID,
	)
	if err != nil {
		return fmt.Errorf("can't update collection: %w", err)
	}
	return nil
}

func (cr *CollectionRepo) generateID() string {
	b := make([]rune, 5)
	for i := range b {
		b[i] = letterRunes[cr.rand.IntN(len(letterRunes))]
	}
	return string(b)
}

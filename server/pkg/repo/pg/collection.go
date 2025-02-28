package pg

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"dishdash.ru/pkg/domain"

	"github.com/Vaniog/go-postgis"
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
            c.id,
            c.name,
            c.description,
            c.avatar,
            c.visible,
            c.updated_at,
            c.created_at,
            JSON_AGG(JSON_BUILD_OBJECT('id', p.id, 'title', p.title, 'shortDescription', p.short_description, 'description', p.description, 'images', p.images, 'location', p.location, 'address', p.address, 'priceAvg', p.price_avg, 'reviewRating', p.review_rating, 'reviewCount', p.review_count, 'source', p.source, 'url', p.url, 'boost', p.boost, 'boost_radius', p.boost_radius)) AS places
            FROM "collection" AS c
			LEFT JOIN "collection_place" AS pc ON c.id = pc.collection_id
			LEFT JOIN "place" AS p ON pc.place_id = p.id
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
            p.source,
            p.url,
            p.boost,
            p.boost_radius,
            p.updated_at
        FROM "place" AS p
        JOIN "collection_place" AS cp ON p.id = cp.place_id
        WHERE cp.collection_id = $1;
    `

	rows, err := cr.db.Query(ctx, getPlacesQuery, collectionID)
	if err != nil {
		return nil, fmt.Errorf("can't fetch places: %w", err)
	}
	defer rows.Close()

	places := make([]*domain.Place, 0)
	for rows.Next() {
		place := new(domain.Place)
		loc := postgis.PointS{}
		var images string
		err := rows.Scan(
			&place.ID,
			&place.Title,
			&place.ShortDescription,
			&place.Description,
			&images,
			&loc,
			&place.Address,
			&place.PriceAvg,
			&place.ReviewRating,
			&place.ReviewCount,
			&place.Source,
			&place.Url,
			&place.Boost,
			&place.BoostRadius,
			&place.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan place: %w", err)
		}
		place.Location = domain.FromPostgis(loc)
		place.Images = strings.Split(images, ",")
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

func scanCollection(s Scanner) (*domain.Collection, error) {
	collection := new(domain.Collection)
	placesJSON := ""

	err := s.Scan(
		&collection.ID,
		&collection.Name,
		&collection.Description,
		&collection.Avatar,
		&collection.Visible,
		&collection.UpdatedAt,
		&collection.CreatedAt,
		&collection.Order,
		&placesJSON,
	)
	if err != nil {
		return nil, err
	}

	var places []map[string]interface{}
	if err := json.Unmarshal([]byte(placesJSON), &places); err != nil {
		return nil, fmt.Errorf("failed to unmarshal places: %w", err)
	}

	collection.Places = make([]*domain.Place, len(places))
	for i, placeData := range places {
		place := &domain.Place{}
		if title, ok := placeData["title"].(string); ok {
			place.Title = title
		}
		if description, ok := placeData["description"].(string); ok {
			place.Description = description
		}
		if images, ok := placeData["images"].(string); ok {
			place.Images = strings.Split(images, ",")
		}
		if location, ok := placeData["location"].(domain.Coordinate); ok {
			place.Location = location
		}
		if address, ok := placeData["address"].(string); ok {
			place.Address = address
		}
		if priceAvg, ok := placeData["priceAvg"].(int); ok {
			place.PriceAvg = priceAvg
		}
		if reviewRating, ok := placeData["reviewRating"].(float64); ok {
			place.ReviewRating = reviewRating
		}
		if reviewCount, ok := placeData["reviewCount"].(int); ok {
			place.ReviewCount = reviewCount
		}
		if updatedAt, ok := placeData["updatedAt"].(string); ok {
			parsedTime, err := time.Parse(time.RFC3339, updatedAt)
			if err != nil {
				return nil, fmt.Errorf("failed to parse time: %w", err)
			}
			place.UpdatedAt = parsedTime
		}
		if source, ok := placeData["source"].(string); ok {
			place.Source = source
		}
		if url, ok := placeData["url"].(string); ok {
			place.Url = &url
		}
		if boost, ok := placeData["boost"].(float64); ok {
			place.Boost = &boost
		}
		if boostRadius, ok := placeData["boostRadius"].(float64); ok {
			place.BoostRadius = &boostRadius
		}

		collection.Places[i] = place
	}

	return collection, nil
}

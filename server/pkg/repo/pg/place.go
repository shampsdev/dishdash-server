package pg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"dishdash.ru/pkg/domain"
	"dishdash.ru/pkg/repo"
	"github.com/Vaniog/go-postgis"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaceRepo struct {
	db *pgxpool.Pool
}

func NewPlaceRepo(db *pgxpool.Pool) *PlaceRepo {
	return &PlaceRepo{db: db}
}

func (pr *PlaceRepo) SavePlace(ctx context.Context, place *domain.Place) (int64, error) {
	const saveQuery = `
	INSERT INTO "place" (
		"title",
		"short_description",
		"description",
		"images",
		"location",
		"address",
		"price_avg",
	    "review_rating",
		"review_count",
		"updated_at",
	    "source",
		"url",
		"boost",
		"boost_radius"
	) VALUES ($1, $2, $3, $4, GeomFromEWKB($5), $6, $7, $8, $9, $10, $11, $12, $13, $14)
	RETURNING "id"
`

	place.UpdatedAt = time.Now().UTC()
	row := pr.db.QueryRow(ctx, saveQuery,
		place.Title,
		place.ShortDescription,
		place.Description,
		strings.Join(place.Images, ","),
		place.Location.ToPostgis(),
		place.Address,
		place.PriceAvg,
		place.ReviewRating,
		place.ReviewCount,
		place.UpdatedAt,
		place.Source,
		place.Url,
		place.Boost,
		place.BoostRadius,
	)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("can't save place: %w", err)
	}
	return id, nil
}

func (pr *PlaceRepo) UpdatePlace(ctx context.Context, place *domain.Place) error {
	const updateQuery = `
	UPDATE "place" SET
	  "title" = $1,
	  "short_description" = $2,
	  "description" = $3,
	  "images" = $4,
	  "location" = GeomFromEWKB($5),
	  "address" = $6,
	  "price_avg" = $7,
	  "review_rating" = $8,
	  "review_count" = $9,
	  "updated_at" = $10,
	  "source" = $11,
	  "url" = $12,
	  "boost" = $13,
	  "boost_radius" = $14
	WHERE "id" = $15
  `
	place.UpdatedAt = time.Now().UTC()
	_, err := pr.db.Exec(ctx, updateQuery,
		place.Title,
		place.ShortDescription,
		place.Description,
		strings.Join(place.Images, ","),
		place.Location.ToPostgis(),
		place.Address,
		place.PriceAvg,
		place.ReviewRating,
		place.ReviewCount,
		place.UpdatedAt,
		place.Source,
		place.Url,
		place.Boost,
		place.BoostRadius,
		place.ID,
	)
	if err != nil {
		return fmt.Errorf("can't update place: %w", err)
	}

	return nil
}

func (pr *PlaceRepo) DeletePlace(ctx context.Context, id int64) error {
	const deleteQuery = `
	DELETE FROM "place"
	WHERE id=$1
`
	_, err := pr.db.Exec(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("can't delete place: %w", err)
	}
	return nil
}

func (pr *PlaceRepo) GetPlaceByID(ctx context.Context, id int64) (*domain.Place, error) {
	const getPlaceQuery = `
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
				JSON_BUILD_OBJECT('id', t.id, 'name', t.name, 'icon', t.icon, 'visible', t.visible, 'order', t.order, 'excluded', t.excluded)
			) FILTER (WHERE t.id IS NOT NULL),
			'[]'
		) AS tags
	FROM "place" AS p
	LEFT JOIN "place_tag" AS pt ON p.id = pt.place_id
	LEFT JOIN "tag" AS t ON pt.tag_id = t.id
	WHERE p.id=$1
	GROUP BY p.id;
`
	row := pr.db.QueryRow(ctx, getPlaceQuery, id)
	place, err := scanPlace(row)
	if err != nil {
		return nil, fmt.Errorf("can't fetch place: %w", err)
	}
	return place, nil
}

func (pr *PlaceRepo) GetPlaceByUrl(ctx context.Context, url string) (*domain.Place, error) {
	const getPlaceQuery = `
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
				JSON_BUILD_OBJECT('id', t.id, 'name', t.name, 'icon', t.icon, 'visible', t.visible, 'order', t.order, 'excluded', t.excluded)
			) FILTER (WHERE t.id IS NOT NULL),
			'[]'
		) AS tags
	FROM "place" AS p
	LEFT JOIN "place_tag" AS pt ON p.id = pt.place_id
	LEFT JOIN "tag" AS t ON pt.tag_id = t.id
	WHERE p.url=$1
	GROUP BY p.id;
`
	row := pr.db.QueryRow(ctx, getPlaceQuery, url)
	place, err := scanPlace(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repo.ErrPlaceNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("can't fetch place: %w", err)
	}
	return place, nil
}

func (pr *PlaceRepo) GetAllPlaces(ctx context.Context) ([]*domain.Place, error) {
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
				JSON_BUILD_OBJECT('id', t.id, 'name', t.name, 'icon', t.icon, 'visible', t.visible, 'order', t.order, 'excluded', t.excluded)
			) FILTER (WHERE t.id IS NOT NULL),
			'[]'
		) AS tags
	FROM "place" AS p
	LEFT JOIN "place_tag" AS pt ON p.id = pt.place_id
	LEFT JOIN "tag" AS t ON pt.tag_id = t.id
	GROUP BY p.id
	ORDER BY p.updated_at DESC;
`
	rows, err := pr.db.Query(ctx, getPlacesQuery)
	if err != nil {
		return nil, fmt.Errorf("can't fetch places: %w", err)
	}
	defer rows.Close()

	places := make([]*domain.Place, 0)
	for rows.Next() {
		place, err := scanPlace(rows)
		if err != nil {
			return nil, fmt.Errorf("can't scan place: %w", err)
		}
		places = append(places, place)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after place scan: %w", err)
	}

	return places, nil
}

func (pr *PlaceRepo) AttachOrderedPlacesToLobby(ctx context.Context, placesIDs []int64, lobbyID string) error {
	if len(placesIDs) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	query := `INSERT INTO place_lobby ("lobby_id", "place_id", "order") VALUES ($1, $2, $3)`
	for i, placeID := range placesIDs {
		batch.Queue(query, lobbyID, placeID, i)
	}

	br := pr.db.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		return fmt.Errorf("could not attach places to lobby: %w", err)
	}
	return nil
}

func (pr *PlaceRepo) DetachPlacesFromLobby(ctx context.Context, placeID string) error {
	query := `DELETE FROM place_lobby WHERE lobby_id=$1`
	_, err := pr.db.Exec(ctx, query, placeID)
	if err != nil {
		return fmt.Errorf("could not detach tags from place: %w", err)
	}
	return nil
}

func (pr *PlaceRepo) GetOrderedPlacesByLobbyID(ctx context.Context, lobbyID string) ([]*domain.Place, error) {
	query := `
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
				JSON_BUILD_OBJECT('id', t.id, 'name', t.name, 'icon', t.icon, 'visible', t.visible, 'order', t.order, 'excluded', t.excluded)
			) FILTER (WHERE t.id IS NOT NULL),
			'[]'
		) AS tags
	FROM "place" AS p
	LEFT JOIN "place_tag" AS pt ON p.id = pt.place_id
	LEFT JOIN "tag" AS t ON pt.tag_id = t.id
	JOIN "place_lobby" AS pl ON p.id = pl.place_id
	WHERE pl.lobby_id = $1
	GROUP BY p.id, pl.order
	ORDER BY pl.order ASC;
	`

	rows, err := pr.db.Query(ctx, query, lobbyID)
	if err != nil {
		return nil, fmt.Errorf("could not get places by lobby ID: %w", err)
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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return places, nil
}

func (pr *PlaceRepo) FilterPlaces(ctx context.Context, filter repo.PlacesFilter) ([]*domain.Place, error) {
	query := `
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
	WHERE p.title ILIKE '%' || $1 || '%'
	`

	args := []any{filter.Search}

	if len(filter.Tags) > 0 {
		query += "\nAND t.name = ANY ($2)"
		args = append(args, filter.Tags)
	}

	query += "\nGROUP BY p.id ORDER BY p.updated_at DESC"

	rows, err := pr.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not filter places: %w", err)
	}
	defer rows.Close()

	places := make([]*domain.Place, 0)
	for rows.Next() {
		place, err := scanPlace(rows)
		if err != nil {
			return nil, fmt.Errorf("could not filter places: %w", err)
		}
		places = append(places, place)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return places, nil
}

type Scanner interface {
	Scan(dest ...any) error
}

func scanPlace(s Scanner) (*domain.Place, error) {
	p := new(domain.Place)
	loc := postgis.PointS{}
	imagesStr := ""
	tagsJSON := ""

	err := s.Scan(
		&p.ID,
		&p.Title,
		&p.ShortDescription,
		&p.Description,
		&imagesStr,
		&loc,
		&p.Address,
		&p.PriceAvg,
		&p.ReviewRating,
		&p.ReviewCount,
		&p.UpdatedAt,
		&p.Source,
		&p.Url,
		&p.Boost,
		&p.BoostRadius,
		&tagsJSON,
	)
	if err != nil {
		return nil, err
	}

	p.Location = domain.FromPostgis(loc)

	p.Images = strings.Split(imagesStr, ",")

	var tags []domain.Tag
	if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	p.Tags = make([]*domain.Tag, len(tags))
	for i := range tags {
		p.Tags[i] = &tags[i]
	}

	return p, nil
}

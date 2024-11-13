package pg

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5"

	"dishdash.ru/cmd/server/config"
	"dishdash.ru/internal/domain"
	"github.com/Vaniog/go-postgis"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	radius             = config.C.Defaults.Radius
	priceAvgLowerDelta = config.C.Defaults.PriceAvgLowerDelta
	priceAvgUpperDelta = config.C.Defaults.PriceAvgUpperDelta
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
		"url"
	) VALUES ($1, $2, $3, $4, GeomFromEWKB($5), $6, $7, $8, $9, $10, $11, $12)
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
	  "url" = $12
	WHERE "id" = $13
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
		place.ID,
	)
	if err != nil {
		return fmt.Errorf("can't update place: %w", err)
	}

	return nil
}

func (pr *PlaceRepo) GetPlaceByID(ctx context.Context, id int64) (*domain.Place, error) {
	const getPlaceQuery = `
	SELECT
		"id",
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
		"url"
	FROM "place"
	WHERE id=$1
`
	row := pr.db.QueryRow(ctx, getPlaceQuery, id)
	place, err := scanPlace(row)
	if err != nil {
		return nil, fmt.Errorf("can't fetch place: %w", err)
	}
	return place, nil
}

func (pr *PlaceRepo) GetAllPlaces(ctx context.Context) ([]*domain.Place, error) {
	const getPlacesQuery = `
	SELECT
		"id",
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
		"url"
	FROM "place"
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

func (pr *PlaceRepo) AttachPlacesToLobby(ctx context.Context, placesIDs []int64, lobbyID string) error {
	if len(placesIDs) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	query := `INSERT INTO place_lobby (lobby_id, place_id) VALUES ($1, $2)`
	for _, placeID := range placesIDs {
		batch.Queue(query, lobbyID, placeID)
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

func (pr *PlaceRepo) GetPlacesByLobbyID(ctx context.Context, lobbyID string) ([]*domain.Place, error) {
	query := `
	SELECT 
	    place.id,
		place.title,
		place.short_description,
		place.description,
		place.images,
		place.location,
		place.address,
		place.price_avg,
		place.review_rating,
		place.review_count,
		place.updated_at,
		place.source,
		place.url
	FROM place
	JOIN place_lobby ON place.id = place_lobby.place_id
	WHERE place_lobby.lobby_id = $1
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

type Scanner interface {
	Scan(dest ...any) error
}

func scanPlace(s Scanner) (*domain.Place, error) {
	p := new(domain.Place)
	loc := postgis.PointS{}
	imagesStr := ""
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
	)
	p.Location = domain.FromPostgis(loc)
	p.Images = strings.Split(imagesStr, ",")
	return p, err
}

func parseTagsToQuery(lobby *domain.Lobby) string {
	query := "HAVING COUNT(DISTINCT CASE WHEN t.name IN (%s) THEN t.name END) > 0"
	var queryTags string
	for _, tag := range lobby.Tags {
		queryTags += fmt.Sprintf("'%s', ", tag.Name)
	}
	queryTags = strings.TrimSuffix(queryTags, ", ")

	return fmt.Sprintf(query, queryTags)
}

func (pr *PlaceRepo) GetPlacesForLobby(ctx context.Context, lobby *domain.Lobby) ([]*domain.Place, error) {
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
			p.url
		FROM place p
		JOIN place_tag pt ON p.id = pt.place_id
		JOIN tag t ON pt.tag_id = t.id
		WHERE ST_DWithin(
				p.location,
				ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
				$3
			)
		  AND p.id IN (
				SELECT pt.place_id 
				FROM place_tag pt
				JOIN tag t ON pt.tag_id = t.id
				GROUP BY pt.place_id
				%s
			)
		  AND p.price_avg > $4
		  AND p.price_avg < $5
		  GROUP BY p.id
	`

	query = fmt.Sprintf(query, parseTagsToQuery(lobby))

	rows, err := pr.db.Query(ctx, query,
		lobby.Location.Lon,
		lobby.Location.Lat,
		radius,
		lobby.PriceAvg-priceAvgLowerDelta,
		lobby.PriceAvg+priceAvgUpperDelta,
	)

	log.Infof("Lobby settings: priceavg: %d - %d", lobby.PriceAvg-priceAvgLowerDelta, lobby.PriceAvg+priceAvgUpperDelta)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var places []*domain.Place

	for rows.Next() {
		place, err := scanPlace(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		places = append(places, place)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	log.Debugf("Total places from database: %d", len(places))
	return places, nil
}

func (pr *PlaceRepo) SaveTwoGisPlace(ctx context.Context, twogisPlace *domain.TwoGisPlace) (int64, error) {
	var existingID int64

	log.Debugf("Checking if place with title '%s' and address '%s' exists.", twogisPlace.Name, twogisPlace.Address)

	err := pr.db.QueryRow(ctx, `
    SELECT id FROM "place"
    WHERE "title" = $1 AND "address" = $2;`, twogisPlace.Name, twogisPlace.Address).Scan(&existingID)
	if err != nil {
		log.Debugf("Error after executing query: %v", err)

		if strings.Contains(err.Error(), "no rows in result set") {
			log.Debug("No rows found, adding new place.")
			place := twogisPlace.ToPlace()
			id, err := pr.SavePlace(ctx, place)
			if err != nil {
				log.WithError(err).Error("Failed to save new place")
				return 0, fmt.Errorf("failed to save place: %w", err)
			}
			log.Debugf("New place saved successfully. ID: %d", id)
			return id, nil
		}

		log.WithError(err).Error("Unexpected error")
		return 0, err
	}

	log.Debugf("Place already exists. ID: %d", existingID)
	return existingID, nil
}

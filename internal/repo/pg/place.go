package pg

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"dishdash.ru/internal/domain"
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
		"updated_at"
	) VALUES ($1, $2, $3, $4, GeomFromEWKB($5), $6, $7, $8, $9, $10)
	RETURNING "id"
`

	place.UpdatedAt = time.Now().UTC()
	row := pr.db.QueryRow(ctx, saveQuery,
		place.Title,
		place.ShortDescription,
		place.Description,
		strings.Join(place.Images, ","),
		postgis.PointS{SRID: 4326, X: place.Location.Lon, Y: place.Location.Lat},
		place.Address,
		place.PriceAvg,
		place.ReviewRating,
		place.ReviewCount,
		place.UpdatedAt,
	)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.Printf("Error saving place: %v\n", err)
		return 0, err
	}
	return id, nil
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
		"updated_at"
	FROM "place"
	WHERE id=$1
`
	row := pr.db.QueryRow(ctx, getPlaceQuery, id)
	place, err := scanPlace(row)
	if err != nil {
		log.Printf("Error fetching places: %v\n", err)
		return nil, err
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
		"updated_at"
	FROM "place"
`
	rows, err := pr.db.Query(ctx, getPlacesQuery)
	if err != nil {
		log.Printf("Error fetching places: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	places := make([]*domain.Place, 0)
	for rows.Next() {
		place, err := scanPlace(rows)
		if err != nil {
			log.Printf("Error scanning place: %v\n", err)
			return nil, err
		}
		places = append(places, place)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after scanning places: %v\n", err)
		return nil, err
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
		return fmt.Errorf("could not attach tags to place: %w", err)
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
		place.updated_at
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
	)
	p.Location = domain.Coordinate{
		Lon: loc.X,
		Lat: loc.Y,
	}
	p.Images = strings.Split(imagesStr, ",")
	return p, err
}
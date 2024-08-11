package pg

import (
	"context"
	"log"
	"time"

	"dishdash.ru/internal/domain"
	"github.com/Vaniog/go-postgis"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApiPlaceRepo struct {
	db *pgxpool.Pool
}

func NewApiPlaceRepo(db *pgxpool.Pool) *ApiPlaceRepo {
	return &ApiPlaceRepo{db: db}
}

func (pr *ApiPlaceRepo) SaveApiPlace(ctx context.Context, place *domain.TwoGisPlace) (int64, error) {
	var exists bool

	err := pr.db.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1
            FROM "place"
            WHERE "title" = $1 AND "address" = $2
        );
    `, place.Name, place.Address).Scan(&exists)
	if err != nil {
		return 0, err
	}

	if exists {
		log.Printf("[INFO] Place with title '%s' and address '%s' already exists", place.Name, place.Address)
		return 0, err
	}

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

	updatedTime := time.Now().UTC()
	row := pr.db.QueryRow(ctx, saveQuery,
		place.Name,
		place.Address,
		place.Address,
		place.PhotoURL,
		postgis.PointS{SRID: 4326, X: place.Lon, Y: place.Lat},
		place.Address,
		place.AveragePrice,
		place.ReviewRating,
		place.ReviewCount,
		updatedTime,
	)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		log.Printf("[ERROR] Error saving place: %v\n", err)
		return 0, err
	}
	return id, nil
}

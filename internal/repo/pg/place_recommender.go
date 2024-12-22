package pg

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"dishdash.ru/internal/domain"
	"github.com/Vaniog/go-postgis"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaceRecommender struct {
	db *pgxpool.Pool
}

func NewPlaceRecommenderRepo(db *pgxpool.Pool) *PlaceRecommender {
	return &PlaceRecommender{db: db}
}

func (pr *PlaceRecommender) RecommendClassic(
	ctx context.Context,
	opts domain.ClassicRecommendationOpts,
	data domain.RecommendData,
) ([]*domain.Place, error) {
	query := `
	WITH filtered_places AS (
		SELECT DISTINCT p.id
		FROM place p
		JOIN place_tag pt ON p.id = pt.place_id
		JOIN tag t ON pt.tag_id = t.id
		WHERE t.name = ANY ($1)
	)
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
		JSON_AGG(
			JSON_BUILD_OBJECT(
				'id', t.id,
				'name', t.name,
				'icon', t.icon,
				'visible', t.visible,
				'order', t.order
			)
		) AS tags
	FROM place p
	JOIN filtered_places fp ON p.id = fp.id
	JOIN place_tag pt ON p.id = pt.place_id
	JOIN tag t ON pt.tag_id = t.id
	GROUP BY p.id
	ORDER BY
		$2 * (ST_Distance(p.location, ST_GeogFromWkb($3)) ^ $4) +
		$5 * (ABS(p.price_avg - $6) ^ $7)
	`

	return pr.queryPlaces(ctx, query,
		data.Tags,
		opts.DistCoeff, data.Location.ToPostgis(), opts.DistPower,
		opts.PriceCoeff, data.PriceAvg, opts.PricePower,
	)
}

func (pr *PlaceRecommender) RecommendPriceBound(
	ctx context.Context,
	opts domain.PriceBoundRecommendationOpts,
	data domain.RecommendData,
) ([]*domain.Place, error) {
	query := `
	WITH filtered_places AS (
		SELECT DISTINCT p.id
		FROM place p
		JOIN place_tag pt ON p.id = pt.place_id
		JOIN tag t ON pt.tag_id = t.id
		WHERE t.name = ANY ($1)
		AND p.price_avg BETWEEN $2 AND $3
	)
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
		JSON_AGG(
			JSON_BUILD_BUILD_OBJECT(
				'id', t.id,
				'name', t.name,
				'icon', t.icon,
				'visible', t.visible,
				'order', t.order
			)
		) AS tags
	FROM place p
	JOIN filtered_places fp ON p.id = fp.id
	JOIN place_tag pt ON p.id = pt.place_id
	JOIN tag t ON pt.tag_id = t.id
	GROUP BY p.id
	ORDER BY
		ST_Distance(p.location, ST_GeogFromWkb($4))
	`

	return pr.queryPlaces(ctx, query,
		data.Tags,
		data.PriceAvg-opts.PriceBound, data.PriceAvg+opts.PriceBound,
		data.Location.ToPostgis(),
	)
}

func (pr *PlaceRecommender) queryPlaces(ctx context.Context, query string, params ...interface{}) ([]*domain.Place, error) {
	rows, err := pr.db.Query(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("error while place recommending from db: %w", err)
	}
	defer rows.Close()

	var places []*domain.Place

	for rows.Next() {
		var place domain.Place
		var tagsJSON []byte
		loc := postgis.PointS{}
		imagesStr := ""

		err := rows.Scan(
			&place.ID, &place.Title, &place.ShortDescription,
			&place.Description, &imagesStr, &loc,
			&place.Address, &place.PriceAvg, &place.ReviewRating,
			&place.ReviewCount, &place.UpdatedAt, &place.Source,
			&place.Url, &tagsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		place.Images = strings.Split(imagesStr, ",")
		place.Location = domain.FromPostgis(loc)

		var tags []*domain.Tag
		if err := json.Unmarshal(tagsJSON, &tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags JSON: %w", err)
		}
		place.Tags = tags

		places = append(places, &place)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return places, nil
}

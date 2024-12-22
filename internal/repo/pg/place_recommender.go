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

func (pr *PlaceRecommender) RecommendPlaces(
	ctx context.Context,
	opts domain.RecommendOpts,
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
		$2 * ST_Distance(p.location, ST_GeogFromWkb($3)) +
		$4 * ABS(p.price_avg - $5)
	`

	rows, err := pr.db.Query(ctx, query,
		data.Tags,
		opts.DistCoeff, data.Location.ToPostgis(),
		opts.PriceCoeff, data.PriceAvg,
	)
	if err != nil {
		return nil, fmt.Errorf("error while place recommending from db: %w", err)
	}
	defer rows.Close()

	var places []*domain.Place

	// TODO: duplicate with repo.pg.place : scanPlace
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

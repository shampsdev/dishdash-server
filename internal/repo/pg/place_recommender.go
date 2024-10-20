package pg

import (
	"context"
	"fmt"

	"dishdash.ru/internal/domain"
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
		p.updated_at
	FROM place p
	JOIN place_tag pt ON p.id = pt.place_id
	JOIN tag t ON pt.tag_id = t.id

	WHERE t.name = ANY ($1)
	
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

	for _, p := range places {
		p.Tags, err = pr.GetTagsByPlaceID(ctx, p.ID)
		if err != nil {
			return nil, err
		}
	}

	return places, nil
}

// TODO: Do in one query
// TODO: duplicate with TagRepo.GetTagsByPlaceID
func (pr *PlaceRecommender) GetTagsByPlaceID(ctx context.Context, placeID int64) ([]*domain.Tag, error) {
	query := `
	SELECT tag.id, tag.name, tag.icon
	FROM tag
	JOIN place_tag ON tag.id = place_tag.tag_id
	WHERE place_tag.place_id = $1
	`

	rows, err := pr.db.Query(ctx, query, placeID)
	if err != nil {
		return nil, fmt.Errorf("could not get tags by place ID: %w", err)
	}
	defer rows.Close()

	tags := make([]*domain.Tag, 0)
	for rows.Next() {
		var tag domain.Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Icon)
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

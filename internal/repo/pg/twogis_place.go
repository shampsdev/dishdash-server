package pg

import (
	"context"
	"fmt"
	"log"

	"dishdash.ru/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ApiPlaceRepo struct {
	db *pgxpool.Pool
}

func NewApiPlaceRepo(db *pgxpool.Pool) *ApiPlaceRepo {
	return &ApiPlaceRepo{db: db}
}

func (apr *ApiPlaceRepo) SaveApiPlace(ctx context.Context, twogisPlace *domain.TwoGisPlace) (int64, error) {
	var exists bool

	err := apr.db.QueryRow(ctx, `
    SELECT EXISTS (
        SELECT 1
        FROM "place"
        WHERE "title" = $1 AND "address" = $2
    );`, twogisPlace.Name, twogisPlace.Address).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("error checking existence of place: %w", err)
	}

	if exists {
		log.Printf("[INFO] Place with title '%s' and address '%s' already exists", twogisPlace.Name, twogisPlace.Address)
		return 0, nil
	}

	place := twogisPlace.ToPlace()

	placeRepo := NewPlaceRepo(apr.db)

	id, err := placeRepo.SavePlace(ctx, place)
	if err != nil {
		log.Printf("Error saving new place: %v\n", err)
		return 0, err
	}

	return id, nil
}

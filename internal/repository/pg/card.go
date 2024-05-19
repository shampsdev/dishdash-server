package pg

import (
	"context"
	"log"

	"dishdash.ru/internal/dto"

	"dishdash.ru/internal/domain"
	"github.com/jackc/pgx/v4"
)

const saveQuery = `
	INSERT INTO "card" (
		"title", 
		"short_description", 
		"description", 
		"image", 
		"location", 
		"address", 
		"type", 
		"price"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING "id"
`

const getCardsQuery = `
	SELECT 
		"id", 
		"title", 
		"short_description", 
		"description", 
		"image", 
		"location", 
		"address", 
		"type", 
		"price"
	FROM "card"
`

type CardRepository struct {
	db *pgx.Conn
}

func NewCardRepository(db *pgx.Conn) *CardRepository {
	return &CardRepository{db: db}
}

func (cr *CardRepository) SaveCard(ctx context.Context, card *domain.Card) error {
	row := cr.db.QueryRow(ctx, saveQuery,
		card.Title,
		card.ShortDescription,
		card.Description,
		card.Image,
		card.Location,
		card.Address,
		card.Type,
		card.Price,
	)

	err := row.Scan(&card.ID)
	if err != nil {
		log.Printf("Error saving card: %v\n", err)
		return err
	}

	return nil
}

func (cr *CardRepository) GetCards(ctx context.Context) ([]*domain.Card, error) {
	rows, err := cr.db.Query(ctx, getCardsQuery)
	if err != nil {
		log.Printf("Error fetching cards: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var cards []*domain.Card
	for rows.Next() {
		var cardDto dto.Card
		err := rows.Scan(
			&cardDto.ID,
			&cardDto.Title,
			&cardDto.ShortDescription,
			&cardDto.Description,
			&cardDto.Image,
			&cardDto.Location,
			&cardDto.Address,
			&cardDto.Type,
			&cardDto.Price,
		)
		if err != nil {
			log.Printf("Error scanning card: %v\n", err)
			return nil, err
		}
		var card *domain.Card
		_ = card.ParseDto(cardDto)
		cards = append(cards, card)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after scanning cards: %v\n", err)
		return nil, err
	}

	return cards, nil
}

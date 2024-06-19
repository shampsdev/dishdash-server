package pg

import (
	"context"
	"log"

	"dishdash.ru/internal/domain"
	"github.com/Vaniog/go-postgis"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CardRepository struct {
	db *pgxpool.Pool
}

func NewCardRepository(db *pgxpool.Pool) *CardRepository {
	return &CardRepository{db: db}
}

func (cr *CardRepository) CreateCard(ctx context.Context, card *domain.Card) (int64, error) {
	const saveQuery = `
	INSERT INTO "card" (
		"title", 
		"short_description", 
		"description", 
		"image", 
		"location", 
		"address", 
		"price"
	) VALUES ($1, $2, $3, $4, GeomFromEWKB($5), $6, $7)
	RETURNING "id"
`

	row := cr.db.QueryRow(ctx, saveQuery,
		card.Title,
		card.ShortDescription,
		card.Description,
		card.Image,
		postgis.PointS{SRID: 4326, X: card.Location.Lon, Y: card.Location.Lat},
		card.Address,
		card.Price,
	)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		log.Printf("Error saving card: %v\n", err)
		return 0, err
	}

	return id, nil
}

func (cr *CardRepository) GetCardByID(ctx context.Context, id int64) (*domain.Card, error) {
	const getCardQuery = `
	SELECT 
		"id", 
		"title", 
		"short_description", 
		"description", 
		"image", 
		"location", 
		"address", 
		"price"
	FROM "card"
	WHERE id=$1
`
	row := cr.db.QueryRow(ctx, getCardQuery, id)
	card, err := scanCard(row)
	if err != nil {
		log.Printf("Error fetching cards: %v\n", err)
		return nil, err
	}
	return card, nil
}

func (cr *CardRepository) GetAllCards(ctx context.Context) ([]*domain.Card, error) {
	const getCardsQuery = `
	SELECT 
		"id", 
		"title", 
		"short_description", 
		"description", 
		"image", 
		"location", 
		"address", 
		"price"
	FROM "card"
`
	rows, err := cr.db.Query(ctx, getCardsQuery)
	if err != nil {
		log.Printf("Error fetching cards: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var cards []*domain.Card
	for rows.Next() {
		card, err := scanCard(rows)
		if err != nil {
			log.Printf("Error scanning card: %v\n", err)
			return nil, err
		}
		cards = append(cards, card)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after scanning cards: %v\n", err)
		return nil, err
	}

	return cards, nil
}

type Scanner interface {
	Scan(dest ...any) error
}

func scanCard(s Scanner) (*domain.Card, error) {
	card := new(domain.Card)
	cardLocation := postgis.PointS{}
	err := s.Scan(
		&card.ID,
		&card.Title,
		&card.ShortDescription,
		&card.Description,
		&card.Image,
		&cardLocation,
		&card.Address,
		&card.Price,
	)
	card.Location = domain.Coordinate{
		Lon: cardLocation.X,
		Lat: cardLocation.Y,
	}
	return card, err
}

package pg

import (
	"context"
	"errors"
	"math/rand/v2"
	"time"

	"dishdash.ru/internal/usecase"

	"github.com/jackc/pgx/v5"

	"github.com/Vaniog/go-postgis"

	"dishdash.ru/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LobbyRepository struct {
	db   *pgxpool.Pool
	rand *rand.Rand
}

func NewLobbyRepository(db *pgxpool.Pool) *LobbyRepository {
	return &LobbyRepository{
		db:   db,
		rand: rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), rand.Uint64())),
	}
}

func (lr *LobbyRepository) CreateLobby(ctx context.Context, lobby *domain.Lobby) (*domain.Lobby, error) {
	tx, err := lr.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const saveLobbyQuery = `
		INSERT INTO "lobby" (
			"id",
			"location",
			"created_at"
		) VALUES ($1, ST_GeogFromWkb($2), $3)
		RETURNING "id", "created_at"
	`
	row := tx.QueryRow(ctx, saveLobbyQuery,
		lr.generateID(),
		postgis.PointS{SRID: 4326, X: lobby.Location.Lat, Y: lobby.Location.Lon},
		time.Now().UTC(),
	)

	err = row.Scan(
		&lobby.ID,
		&lobby.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	const saveLobbySettingsQuery = `
		INSERT INTO "lobbysettings" (
			"lobby_id",
			"price_min",
			"price_max",
			"max_distance"
		) VALUES ($1, $2, $3, $4)
		RETURNING "id"
	`
	settingsRow := tx.QueryRow(ctx, saveLobbySettingsQuery,
		lobby.ID,
		0.0,
		1000000.0,
		1000000.0,
	)

	var settingsID int
	err = settingsRow.Scan(&settingsID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return lobby, nil
}

func (lr *LobbyRepository) NearestLobby(ctx context.Context, location domain.Coordinate) (*domain.Lobby, float64, error) {
	const getQuery = `
	SELECT lobby.id, lobby.created_at, lobby.location, ST_Distance(lobby.location, ST_GeogFromWkb($1)) as dist
    FROM lobby
    WHERE ST_Distance(lobby.location, ST_GeogFromWkb($1), true) = (
    	SELECT MIN (ST_Distance(lobby.location, ST_GeogFromWkb($1))) 
    	FROM lobby
  	);
`
	row := lr.db.QueryRow(ctx, getQuery,
		postgis.PointS{SRID: 4326, X: location.Lat, Y: location.Lon},
	)

	lobby := new(domain.Lobby)
	dist := 0.0
	var loc postgis.PointS
	err := row.Scan(
		&lobby.ID,
		&lobby.CreatedAt,
		&loc,
		&dist,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, 0, usecase.ErrLobbyNotFound
	}
	lobby.Location = domain.Coordinate{Lat: loc.X, Lon: loc.Y}
	if err != nil {
		return nil, 0, err
	}
	return lobby, dist, nil
}

func (lr *LobbyRepository) DeleteLobbyByID(ctx context.Context, lobbyID string) error {
	const deleteQuery = `
	WITH deleted as (
		DELETE FROM lobby
		WHERE id = $1
		RETURNING *
	) 
	SELECT count(*) FROM deleted
`
	row := lr.db.QueryRow(ctx, deleteQuery, lobbyID)
	amount := 0
	err := row.Scan(&amount)
	if err != nil {
		return err
	}
	if amount == 0 {
		return errors.New("nothing to delete")
	}
	return nil
}

var letterRunes = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (lr *LobbyRepository) generateID() string {
	b := make([]rune, 5)
	for i := range b {
		b[i] = letterRunes[lr.rand.IntN(len(letterRunes))]
	}
	return string(b)
}

func (lr *LobbyRepository) GetLobbyByID(ctx context.Context, id string) (*domain.Lobby, error) {
	const getQuery = `
	SELECT id, created_at, location
	FROM lobby
	WHERE id = $1
`
	row := lr.db.QueryRow(ctx, getQuery, id)

	var lobby domain.Lobby
	var location postgis.PointS

	err := row.Scan(&lobby.ID, &lobby.CreatedAt, &location)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, usecase.ErrLobbyNotFound
	}
	if err != nil {
		return nil, err
	}
	lobby.Location = domain.Coordinate{Lat: location.X, Lon: location.Y}

	lobbySettings, err := lr.getLobbySettings(ctx, lobby.ID)
	if err != nil {
		return nil, err
	}
	lobby.LobbySettings = lobbySettings

	cards, err := lr.getLobbyCards(ctx, lobby.ID)
	if err != nil {
		return nil, err
	}
	lobby.Cards = cards

	matches, err := lr.getMatches(ctx, lobby.ID)
	if err != nil {
		return nil, err
	}
	lobby.Matches = matches

	finalVotes, err := lr.getFinalVotes(ctx, lobby.ID)
	if err != nil {
		return nil, err
	}
	lobby.FinalVotes = finalVotes

	swipes, err := lr.getSwipes(ctx, lobby.ID)
	if err != nil {
		return nil, err
	}
	lobby.Swipes = swipes

	return &lobby, nil
}

func (lr *LobbyRepository) getLobbySettings(ctx context.Context, lobbyID string) (*domain.LobbySettings, error) {
	const query = `
		SELECT id, lobby_id, price_min, price_max, max_distance
		FROM lobbysettings
		WHERE lobby_id = $1
	`
	row := lr.db.QueryRow(ctx, query, lobbyID)

	var lobbySettings domain.LobbySettings
	if err := row.Scan(&lobbySettings.ID, &lobbySettings.LobbyID, &lobbySettings.PriceMin, &lobbySettings.PriceMax, &lobbySettings.MaxDistance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &lobbySettings, nil
}

func (lr *LobbyRepository) getLobbyCards(ctx context.Context, lobbyID string) ([]*domain.Card, error) {
	const query = `
		SELECT c.id, c.title, c.short_description, c.description, c.image, c.location, c.address, c.price_min, c.price_max
		FROM lobby_card lc
		JOIN card c ON lc.card_id = c.id
		WHERE lc.lobby_id = $1
	`
	rows, err := lr.db.Query(ctx, query, lobbyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []*domain.Card
	for rows.Next() {
		var card domain.Card
		var location postgis.PointS
		if err := rows.Scan(&card.ID, &card.Title, &card.ShortDescription, &card.Description, &card.Image, &location, &card.Address, &card.PriceMin, &card.PriceMax); err != nil {
			return nil, err
		}
		card.Location = domain.Coordinate{Lat: location.X, Lon: location.Y}
		cards = append(cards, &card)
	}

	return cards, nil
}

func (lr *LobbyRepository) getMatches(ctx context.Context, lobbyID string) ([]*domain.Match, error) {
	const query = `
		SELECT m.id, m.lobby_id, m.card_id
		FROM match m
		WHERE m.lobby_id = $1
	`
	rows, err := lr.db.Query(ctx, query, lobbyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []*domain.Match
	for rows.Next() {
		var match domain.Match
		if err := rows.Scan(&match.ID, &match.LobbyID, &match.CardID); err != nil {
			return nil, err
		}
		matches = append(matches, &match)
	}

	return matches, nil
}

func (lr *LobbyRepository) getFinalVotes(ctx context.Context, lobbyID string) ([]*domain.FinalVote, error) {
	const query = `
		SELECT fv.id, fv.lobby_id, fv.card_id, fv.user_id
		FROM final_vote fv
		WHERE fv.lobby_id = $1
	`
	rows, err := lr.db.Query(ctx, query, lobbyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var finalVotes []*domain.FinalVote
	for rows.Next() {
		var finalVote domain.FinalVote
		if err := rows.Scan(&finalVote.ID, &finalVote.LobbyID, &finalVote.CardID, &finalVote.UserID); err != nil {
			return nil, err
		}
		finalVotes = append(finalVotes, &finalVote)
	}

	return finalVotes, nil
}

func (lr *LobbyRepository) getSwipes(ctx context.Context, lobbyID string) ([]*domain.Swipe, error) {
	const query = `
		SELECT sw.lobby_id, sw.card_id, sw.user_id, sw.type
		FROM swipe sw
		WHERE sw.lobby_id = $1
	`
	rows, err := lr.db.Query(ctx, query, lobbyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var swipes []*domain.Swipe
	for rows.Next() {
		var swipe domain.Swipe
		if err := rows.Scan(
			&swipe.LobbyID,
			&swipe.CardID,
			&swipe.UserID,
			&swipe.Type,
		); err != nil {
			return nil, err
		}
		swipes = append(swipes, &swipe)
	}

	return swipes, nil
}

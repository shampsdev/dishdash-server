package pg

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"

	"github.com/Vaniog/go-postgis"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LobbyRepo struct {
	db   *pgxpool.Pool
	rand *rand.Rand
}

func NewLobbyRepo(db *pgxpool.Pool) *LobbyRepo {
	return &LobbyRepo{
		db:   db,
		rand: rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), rand.Uint64())),
	}
}

func (lr *LobbyRepo) SaveLobby(ctx context.Context, lobby *domain.Lobby) (string, error) {
	const saveQuery = `
		INSERT INTO "lobby" (id, state, price_avg, location, created_at)
		VALUES ($1, $2, $3, ST_GeogFromWkb($4), $5)
`

	lobby.ID = lr.generateID()
	lobby.CreatedAt = time.Now().UTC()
	_, err := lr.db.Exec(ctx, saveQuery,
		lobby.ID,
		lobby.State,
		lobby.PriceAvg,
		lobby.Location.ToPostgis(),
		lobby.CreatedAt,
	)
	if err != nil {
		return "", fmt.Errorf("can't save lobby: %w", err)
	}
	return lobby.ID, nil
}

func (lr *LobbyRepo) DeleteLobbyByID(ctx context.Context, id string) error {
	const deleteQuery = `
		DELETE FROM "lobby" WHERE id = $1
`
	_, err := lr.db.Exec(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("can't delete lobby: %w", err)
	}
	return nil
}

func (lr *LobbyRepo) GetLobbyByID(ctx context.Context, id string) (*domain.Lobby, error) {
	const getQuery = `
		SELECT id, state, price_avg, location, created_at
		FROM "lobby" WHERE id = $1
`
	row := lr.db.QueryRow(ctx, getQuery, id)

	lobby := &domain.Lobby{}
	var loc postgis.PointS
	err := row.Scan(
		&lobby.ID,
		&lobby.State,
		&lobby.PriceAvg,
		&loc,
		&lobby.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("can't get lobby: %w", err)
	}
	lobby.Location = domain.FromPostgis(loc)
	return lobby, nil
}

// NOTE
// This function finds the nearest active lobby
// that was created no longer than 2 min ago.
func (lr *LobbyRepo) NearestActiveLobbyID(ctx context.Context, loc domain.Coordinate) (string, float64, error) {
	const getQuery = `
	SELECT lobby.id, ST_Distance(lobby.location, ST_GeogFromWkb($1)) as dist
	FROM lobby
	WHERE ST_Distance(lobby.location, ST_GeogFromWkb($1), true) = (
	    SELECT MIN(ST_Distance(lobby.location, ST_GeogFromWkb($1)))
	    FROM lobby
	    WHERE lobby.state = 'lobby'
	)
	AND lobby.created_at >= NOW() - INTERVAL '30 seconds';
`
	row := lr.db.QueryRow(ctx, getQuery,
		loc.ToPostgis(),
	)

	id := ""
	dist := 0.0
	err := row.Scan(
		&id,
		&dist,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", 0, repo.ErrLobbyNotFound
	}
	return id, dist, nil
}

func (lr *LobbyRepo) UpdateLobby(ctx context.Context, lobby *domain.Lobby) error {
	const query = `
		UPDATE lobby SET price_avg = $1, location = ST_GeogFromWkb($2)
		WHERE id = $3
`
	_, err := lr.db.Exec(ctx, query,
		lobby.PriceAvg,
		lobby.Location.ToPostgis(),
		lobby.ID,
	)
	if err != nil {
		return fmt.Errorf("can't update lobby: %w", err)
	}
	return nil
}

func (lr *LobbyRepo) SetLobbyState(ctx context.Context, lobbyID string, state domain.LobbyState) error {
	const query = `
		UPDATE lobby SET state = $1
		WHERE id = $2
`
	_, err := lr.db.Exec(ctx, query, state, lobbyID)
	if err != nil {
		return fmt.Errorf("can't update lobby state: %w", err)
	}
	return nil
}

var letterRunes = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (lr *LobbyRepo) generateID() string {
	b := make([]rune, 5)
	for i := range b {
		b[i] = letterRunes[lr.rand.IntN(len(letterRunes))]
	}
	return string(b)
}

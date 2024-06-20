package pg

import (
	"context"
	"errors"
	"github.com/Vaniog/go-postgis"
	"math/rand/v2"
	"time"

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

func (lr LobbyRepository) CreateLobby(ctx context.Context, lobby *domain.Lobby) (*domain.Lobby, error) {
	const saveQuery = `
	INSERT INTO "lobby" (
		"id",
		"location",
		"created_at"
	) VALUES ($1, ST_GeogFromWkb($2), $3)
	RETURNING "id", "created_at"
`
	row := lr.db.QueryRow(ctx, saveQuery,
		lr.generateID(),
		postgis.PointS{SRID: 4326, X: lobby.Location.Lat, Y: lobby.Location.Lon},
		time.Now().UTC(),
	)

	err := row.Scan(
		&lobby.ID,
		&lobby.CreatedAt,
	)

	return lobby, err
}

func (lr LobbyRepository) NearestLobby(ctx context.Context, location domain.Coordinate) (*domain.Lobby, float64, error) {
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
	lobby.Location = domain.Coordinate{Lat: loc.X, Lon: loc.Y}
	if err != nil {
		return nil, 0, err
	}
	return lobby, dist, nil
}

func (lr LobbyRepository) DeleteLobbyByID(ctx context.Context, lobbyID string) error {
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

func (lr LobbyRepository) generateID() string {
	b := make([]rune, 5)
	for i := range b {
		b[i] = letterRunes[lr.rand.IntN(len(letterRunes))]
	}
	return string(b)
}

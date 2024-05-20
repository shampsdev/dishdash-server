package pg

import (
	"context"
	"math/rand"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/dto"
	"github.com/jackc/pgx/v4"
)

type LobbyRepository struct {
	db *pgx.Conn
}

func NewLobbyRepository(db *pgx.Conn) *LobbyRepository {
	return &LobbyRepository{db: db}
}

func (r *LobbyRepository) GetLobbyByID(ctx context.Context, id int64) (*domain.Lobby, error) {
	query := `SELECT id, location FROM lobby WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)

	var lobbyDto dto.Lobby
	if err := row.Scan(&lobbyDto.ID, &lobbyDto.Location); err != nil {
		return nil, err
	}

	lobby := new(domain.Lobby)
	return lobby, lobby.ParseDto(lobbyDto)
}

func randID() int {
	minID := 100000
	maxID := 999999
	return rand.Intn(maxID-minID+1) + minID
}

func (r *LobbyRepository) SaveLobby(ctx context.Context, lobby *domain.Lobby) error {
	id := randID()
	query := `INSERT INTO lobby (id, location) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(ctx, query, id, domain.Point2String(lobby.Location)).Scan(&lobby.ID)
	return err
}

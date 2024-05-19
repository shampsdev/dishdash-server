package pg

import (
	"context"

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

	var lobby *domain.Lobby
	return lobby, lobby.ParseDto(lobbyDto)
}

func (r *LobbyRepository) SaveLobby(ctx context.Context, lobby *domain.Lobby) error {
	query := `INSERT INTO lobby (location) VALUES ($1) RETURNING id`
	err := r.db.QueryRow(ctx, query, domain.Point2String(lobby.Location)).Scan(&lobby.ID)
	return err
}

package pg

import (
	"context"
	"fmt"
	"time"

	"dishdash.ru/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) SaveUser(ctx context.Context, user *domain.User) (string, error) {
	const query = `
		INSERT INTO "user" (name, avatar, telegram, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (telegram) DO UPDATE 
		SET name = EXCLUDED.name, avatar = EXCLUDED.avatar, updated_at = EXCLUDED.updated_at
		RETURNING id
	`
	var id string
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now
	err := ur.db.QueryRow(ctx, query,
		user.Name,
		user.Avatar,
		user.Telegram,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("could not insert or update user: %w", err)
	}
	return id, nil
}

func (ur *UserRepo) SaveUserWithID(ctx context.Context, user *domain.User, id string) error {
	const query = `
		INSERT INTO "user" (id, name, avatar, telegram, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING
	`
	now := time.Now().UTC()
	_, err := ur.db.Exec(ctx, query,
		id,
		user.Name,
		user.Avatar,
		user.Telegram,
		user.CreatedAt,
		now,
	)
	if err != nil {
		return fmt.Errorf("could not insert user: %w", err)
	}
	return nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	const query = `
        UPDATE "user"
        SET name = $2, avatar = $3, telegram = $4, updated_at = $5
        WHERE id = $1
	`
	now := time.Now().UTC()
	commandTag, err := ur.db.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Avatar,
		user.Telegram,
		now,
	)
	if err != nil {
		return user, fmt.Errorf("could not update user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return user, fmt.Errorf("no rows affected, user not found")
	}

	return user, nil
}

func (ur *UserRepo) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	const query = `
		SELECT id, name, avatar, telegram, created_at, updated_at
		FROM "user"
		WHERE id = $1
	`
	user := new(domain.User)
	err := ur.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Avatar,
		&user.Telegram,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get user by id: %w", err)
	}
	return user, nil
}

func (ur *UserRepo) GetUserByTelegram(ctx context.Context, telegram *int64) (*domain.User, error) {
	const query = `
		SELECT id, name, avatar, telegram, created_at, updated_at
		FROM "user"
		WHERE telegram = $1
	`
	user := new(domain.User)
	err := ur.db.QueryRow(ctx, query, telegram).Scan(
		&user.ID,
		&user.Name,
		&user.Avatar,
		&user.Telegram,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get user by id: %w", err)
	}
	return user, nil
}

func (ur *UserRepo) AttachUsersToLobby(ctx context.Context, userIDs []string, lobbyID string) error {
	if len(userIDs) == 0 {
		return nil
	}
	batch := &pgx.Batch{}

	query := `INSERT INTO lobby_user (user_id, lobby_id) VALUES ($1, $2)`
	for _, userID := range userIDs {
		batch.Queue(query, userID, lobbyID)
	}

	br := ur.db.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		return fmt.Errorf("could not attach users to lobby: %w", err)
	}
	return nil
}

func (ur *UserRepo) DetachUsersFromLobby(ctx context.Context, lobbyID string) error {
	query := `DELETE FROM lobby_user WHERE lobby_id = $1`
	_, err := ur.db.Exec(ctx, query, lobbyID)
	if err != nil {
		return fmt.Errorf("could not detach tags from lobby: %w", err)
	}
	return nil
}

func (ur *UserRepo) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	const query = `
		SELECT id, name, avatar, telegram, created_at, updated_at
		FROM "user"
	`

	rows, err := ur.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get all users: %w", err)
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		user := new(domain.User)
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Avatar,
			&user.Telegram,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepo) GetUsersByLobbyID(ctx context.Context, lobbyID string) ([]*domain.User, error) {
	const query = `
		SELECT id, name, avatar, telegram, created_at, updated_at
		FROM "user"
		JOIN lobby_user ON "user".id = lobby_user.user_id
		WHERE lobby_user.lobby_id = $1
	`
	users := make([]*domain.User, 0)
	rows, err := ur.db.Query(ctx, query, lobbyID)
	if err != nil {
		return nil, fmt.Errorf("could not get users by lobby id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := new(domain.User)
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Avatar,
			&user.Telegram,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not get users by lobby id: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not get users by lobby id: %w", err)
	}
	return users, nil
}

package pg

import (
	"context"
	"fmt"
	"time"

	"dishdash.ru/internal/domain"

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
		INSERT INTO "user" (name, avatar, telegram, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	user.CreatedAt = time.Now().UTC()
	err := ur.db.QueryRow(ctx, query,
		user.Name,
		user.Avatar,
		user.Telegram,
		user.CreatedAt,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("could not insert user: %w", err)
	}
	return id, nil
}

func (ur *UserRepo) SaveUserWithID(ctx context.Context, user *domain.User, id string) error {
	const query = `
		INSERT INTO "user" (id, name, avatar, telegram, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := ur.db.Exec(ctx, query,
		id,
		user.Name,
		user.Avatar,
		user.Telegram,
		user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("could not insert user: %w", err)
	}
	return nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	const query = `
        UPDATE "user"
        SET name = $2, avatar = $3, telegram = $4
        WHERE id = $1
	`

	commandTag, err := ur.db.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Avatar,
		user.Telegram,
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
		SELECT id, name, avatar, telegram, created_at
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
	)
	if err != nil {
		return nil, fmt.Errorf("could not get user by id: %w", err)
	}
	return user, nil
}

func (ur *UserRepo) AttachUserToLobby(ctx context.Context, userID, lobbyID string) error {
	const query = `
		INSERT INTO lobby_user (lobby_id, user_id)
		VALUES ($1, $2)
`
	_, err := ur.db.Exec(ctx, query, lobbyID, userID)
	if err != nil {
		return fmt.Errorf("could not attach user to lobby: %w", err)
	}
	return nil
}

func (ur *UserRepo) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	const query = `
		SELECT id, name, avatar, telegram, created_at
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
		SELECT id, name, avatar, telegram, created_at
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

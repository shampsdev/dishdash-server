package pg

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"dishdash.ru/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db   *pgxpool.Pool
	rand *rand.Rand
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db:   db,
		rand: rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), rand.Uint64())),
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	const query = `
		INSERT INTO "user" (id, name, avatar, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id string
	user.CreatedAt = time.Now().UTC()
	err := ur.db.QueryRow(ctx, query, ur.generateID(), user.Name, user.Avatar, user.CreatedAt).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("could not insert user: %w", err)
	}
	return id, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (string, error) {
	const query = `
        UPDATE "user"
        SET name = $2, avatar = $3
        WHERE id = $1
    `

	commandTag, err := ur.db.Exec(ctx, query, user.ID, user.Name, user.Avatar)
	if err != nil {
		return "", fmt.Errorf("could not update user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return "", fmt.Errorf("no rows affected, user not found")
	}

	return user.ID, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	const query = `
	SELECT id, name, avatar, created_at
	FROM "user"
	WHERE id = $1
`
	user := new(domain.User)
	err := ur.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.Avatar, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not get user by id: %w", err)
	}
	return user, nil
}

func (ur *UserRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	const query = `
	SELECT id, name, avatar, created_at
	FROM "user"
`
	rows, err := ur.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := new(domain.User)
		err := rows.Scan(&user.ID, &user.Name, &user.Avatar, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("could not scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (ur *UserRepository) generateID() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = letterRunes[ur.rand.IntN(len(letterRunes))]
	}
	return string(b)
}

package user

import (
	"time"

	"dishdash.ru/internal/domain"
)

type userOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"createdAt"`
}

func userToOutput(u *domain.User) userOutput {
	return userOutput{
		ID:        u.ID,
		Name:      u.Name,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
	}
}

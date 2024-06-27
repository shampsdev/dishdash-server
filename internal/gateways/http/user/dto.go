package user

import (
	"dishdash.ru/internal/domain"
)

type userOutput struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func userToOutput(u *domain.User) userOutput {
	return userOutput{
		ID:     u.ID,
		Name:   u.Name,
		Avatar: u.Avatar,
	}
}

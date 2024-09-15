package usecase

import (
	"context"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"
)

type UserUseCase struct {
	uRepo repo.User
}

func NewUserUseCase(uRepo repo.User) *UserUseCase {
	return &UserUseCase{uRepo: uRepo}
}

func (u UserUseCase) SaveUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	id, err := u.uRepo.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, err
}

func (u UserUseCase) SaveUserWithID(ctx context.Context, user *domain.User, id string) error {
	return u.uRepo.SaveUserWithID(ctx, user, id)
}

func (u UserUseCase) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return u.uRepo.UpdateUser(ctx, user)
}

func (u UserUseCase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return u.uRepo.GetUserByID(ctx, id)
}

func (u UserUseCase) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return u.uRepo.GetAllUsers(ctx)
}

package usecase

import (
	"context"
	"time"

	"dishdash.ru/internal/domain"
	"dishdash.ru/internal/repo"
)

type UserUseCase struct {
	userRepo repo.User
}

func NewUserUseCase(userRepo repo.User) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (u *UserUseCase) CreateUser(ctx context.Context, userInput UserInput) (*domain.User, error) {
	user := &domain.User{
		Name:      userInput.Name,
		Avatar:    userInput.Avatar,
		CreatedAt: time.Now(),
	}
	id, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

func (u *UserUseCase) UpdateUser(ctx context.Context, userInput UserInputExtended) (*domain.User, error) {
	user := &domain.User{
		ID:     userInput.ID,
		Name:   userInput.Name,
		Avatar: userInput.Avatar,
	}
	id, err := u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return user, nil
}

func (u *UserUseCase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}

func (u *UserUseCase) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return u.userRepo.GetAllUsers(ctx)
}

package usecase

import (
	"context"
	"grovo/internal/common"
	"grovo/internal/entity"
)

type (
	UserUsecase struct {
		repo UserRepository
	}
	UserRepository interface {
		Login(ctx context.Context, username string) (*entity.User, error)
		Register(ctx context.Context, username, password string) (*entity.User, error)
		IsUnique(ctx context.Context, username string) (bool, error)
	}
)

func (u *UserUsecase) Login(ctx context.Context, username string) (*entity.User, error) {
	return u.repo.Login(ctx, username)
}

func (u *UserUsecase) Register(ctx context.Context, username string, password string) (*entity.User, error) {
	unique, err := u.repo.IsUnique(ctx, username)
	if err != nil {
		return nil, err
	}
	if !unique {
		return nil, common.ErrUsernameTaken
	}
	return u.repo.Register(ctx, username, password)
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

package usecase

import (
	"cat-test/internal/domain"
	"context"
)

type userUsecase struct {
	repository domain.UserRepository
}

func (u userUsecase) List(ctx context.Context) ([]domain.User, error) {
	return u.repository.List(ctx)
}

func (u userUsecase) Get(ctx context.Context, ID int64) (domain.User, error) {
	return u.repository.Get(ctx, ID)
}

func (u userUsecase) Create(ctx context.Context, user domain.User) (int64, error) {
	return u.repository.Create(ctx, user)
}

func (u userUsecase) Update(ctx context.Context, user domain.User) error {
	return u.repository.Update(ctx, user)
}

func (u userUsecase) Delete(ctx context.Context, ID int64) error {
	return u.repository.Delete(ctx, ID)
}

func NewUserUsecase(repository domain.UserRepository) domain.UserUsecase {
	return userUsecase{repository: repository}
}

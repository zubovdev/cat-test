package usecase

import (
	"cat-test/internal/domain"
	"cat-test/internal/repository"
)

type Usecases struct {
	Auth domain.AuthUsecase
	User domain.UserUsecase
	Task domain.TaskUsecase
}

func NewUsecases(repositories repository.Repositories) Usecases {
	return Usecases{
		Auth: NewAuthUsecase(repositories.Auth, repositories.User),
		User: NewUserUsecase(repositories.User),
		Task: NewTaskUsecase(repositories.Task),
	}
}

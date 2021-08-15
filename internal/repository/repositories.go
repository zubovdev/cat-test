package repository

import (
	"cat-test/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	Auth domain.AuthRepository
	User domain.UserRepository
	Task domain.TaskRepository
}

func NewRepositories(db *sqlx.DB) Repositories {
	return Repositories{
		Auth: NewAuthRepository(db),
		User: NewUserRepository(db),
		Task: NewTaskRepository(db),
	}
}

package domain

import (
	"context"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, userID int64) error
	Authenticate(ctx context.Context, token string) (User, error)
	HashPassword(password string) (string, error)
}

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByToken(ctx context.Context, token string) (User, error)
	CreateToken(ctx context.Context, token string, userID int64) error
	DestroyToken(ctx context.Context, userID int64) error
}

type AuthValidator interface {
	EmailIsUnique(email interface{}) error
}

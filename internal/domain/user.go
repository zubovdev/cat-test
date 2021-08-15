package domain

import (
	"context"
)

const (
	UserType = iota + 1
	UserTypeAdmin
)

type User struct {
	ID           int64   `db:"id" json:"id"`
	Email        string  `db:"email" json:"email"`
	FirstName    *string `db:"first_name" json:"firstName"`
	LastName     *string `db:"last_name" json:"lastName"`
	Type         int     `db:"type" json:"type"`
	PasswordHash string  `db:"password_hash" json:"-"`
	AuthToken    *string `db:"auth_token" json:"-"`
}

type UserUsecase interface {
	List(ctx context.Context) ([]User, error)
	Get(ctx context.Context, ID int64) (User, error)
	Create(ctx context.Context, user User) (int64, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, ID int64) error
}

type UserRepository interface {
	List(ctx context.Context) ([]User, error)
	Get(ctx context.Context, ID int64) (User, error)
	Create(ctx context.Context, user User) (int64, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, ID int64) error
}

type UserValidator interface {
	EmailIsUnique(email interface{}) error
}

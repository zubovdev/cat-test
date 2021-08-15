package repository

import (
	"cat-test/internal/domain"
	"context"
	"github.com/jmoiron/sqlx"
)

type authRepository struct {
	db *sqlx.DB
}

func (a authRepository) GetUserByEmail(ctx context.Context, email string) (user domain.User, err error) {
	sql := `SELECT * FROM "user" WHERE email = $1 LIMIT 1`
	return user, a.db.GetContext(ctx, &user, sql, email)
}

func (a authRepository) GetUserByToken(ctx context.Context, token string) (user domain.User, err error) {
	sql := `SELECT * FROM "user" WHERE auth_token = $1 LIMIT 1`
	return user, a.db.GetContext(ctx, &user, sql, token)
}

func (a authRepository) CreateToken(ctx context.Context, token string, userID int64) error {
	sql := `UPDATE "user" SET auth_token = $1 WHERE id = $2`
	_, err := a.db.ExecContext(ctx, sql, token, userID)
	return err
}

func (a authRepository) DestroyToken(ctx context.Context, userID int64) error {
	sql := `UPDATE "user" SET auth_token = null WHERE id = $1`
	_, err := a.db.ExecContext(ctx, sql, userID)
	return err
}

func NewAuthRepository(db *sqlx.DB) domain.AuthRepository {
	return authRepository{db: db}
}

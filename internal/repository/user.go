package repository

import (
	"cat-test/internal/domain"
	"context"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func (u userRepository) List(ctx context.Context) (users []domain.User, err error) {
	sql := `SELECT * FROM "user" ORDER BY id`
	return users, u.db.SelectContext(ctx, &users, sql)
}

func (u userRepository) Get(ctx context.Context, ID int64) (user domain.User, err error) {
	sql := `SELECT * FROM "user" WHERE "user".id = $1 LIMIT 1`
	return user, u.db.GetContext(ctx, &user, sql, ID)
}

func (u userRepository) Create(ctx context.Context, user domain.User) (userID int64, err error) {
	sql := `INSERT INTO "user"(email, first_name, last_name, type, password_hash) 
				VALUES(:email, :first_name, :last_name, :type, :password_hash) RETURNING id`
	rows, err := u.db.NamedQueryContext(ctx, sql, user)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()

	return userID, rows.Scan(&userID)
}

func (u userRepository) Update(ctx context.Context, user domain.User) error {
	sql := `UPDATE "user" SET email = :email, first_name = :first_name, 
                  last_name = :last_name, type = :type WHERE id = :id`
	_, err := u.db.NamedExecContext(ctx, sql, user)
	return err
}

func (u userRepository) Delete(ctx context.Context, ID int64) error {
	sql := `DELETE FROM "user" WHERE "user".id = $1`
	_, err := u.db.ExecContext(ctx, sql, ID)
	return err
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return userRepository{db: db}
}

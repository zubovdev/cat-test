package validator

import (
	"cat-test/internal/domain"
	"errors"
	"github.com/jmoiron/sqlx"
)

type userValidator struct {
	db *sqlx.DB
}

func (u userValidator) EmailIsUnique(email interface{}) error {
	var count int64
	if err := u.db.Get(&count, `SELECT count(*) FROM "user" WHERE email = $1 LIMIT 1`, email); err != nil {
		return err
	} else if count != 0 {
		return errors.New("user with this email already exist")
	}

	return nil
}

func NewUserValidator(db *sqlx.DB) domain.UserValidator {
	return userValidator{db: db}
}

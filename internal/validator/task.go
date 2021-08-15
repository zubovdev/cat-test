package validator

import (
	"cat-test/internal/domain"
	"errors"
	"github.com/jmoiron/sqlx"
)

type taskValidator struct {
	db *sqlx.DB
}

func (u taskValidator) ValidateUserExist(userID interface{}) error {
	var count int64
	if err := u.db.Get(&count, `SELECT count(1) FROM "user" WHERE id = $1`, userID); err != nil {
		return err
	} else if count == 0 {
		return errors.New("user does not exist")
	}

	return nil
}

func NewTaskValidator(db *sqlx.DB) domain.TaskValidator {
	return taskValidator{db: db}
}

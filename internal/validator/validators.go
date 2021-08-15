package validator

import (
	"cat-test/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Validators struct {
	//Auth domain.AuthValidator
	User domain.UserValidator
	Task domain.TaskValidator
}

func NewValidators(db *sqlx.DB) Validators {
	return Validators{
		//Auth: NewAuthValidator(db),
		User: NewUserValidator(db),
		Task: NewTaskValidator(db),
	}
}

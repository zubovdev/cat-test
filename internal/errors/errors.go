package errors

import "errors"

var (
	AuthInvalidEmail           = errors.New("invalid email")
	AuthInvalidPassword        = errors.New("invalid password")
	AuthInvalidEmailOrPassword = errors.New("invalid email or password")
	TaskCannotBeAssigned       = errors.New("you cannot assign this task")
)

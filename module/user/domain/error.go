package domain

import "errors"

var (
	ErrEmailHasExisted      = errors.New("email has been existed")
	ErrInvalidEmailPassword = errors.New("invalid email and password")
)

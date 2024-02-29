package domain

import "errors"

var (
	ErrInvalidUser               = errors.New("Invalid username or password")
	ErrInvalidDataonRegistartion = errors.New("Invalid username, email or password")
	ErrUserAlreadyExist          = errors.New("User already exist")
	ErrSessionNotFound           = errors.New("session not found")
)

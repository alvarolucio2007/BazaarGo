package domain

import "errors"

var (
	//User
	ErrInvalidCPF         = errors.New("invalid CPF")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("student not found")
)

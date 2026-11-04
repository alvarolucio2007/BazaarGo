package domain

import "errors"

var (
	//Generic
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("already exists")

	//User
	ErrInvalidCPF         = errors.New("invalid CPF")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("student not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrCPFAlreadyExists   = errors.New("CPF already exists")
)

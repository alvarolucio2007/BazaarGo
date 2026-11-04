package domain

import (
	"net/mail"
	"regexp"
	"time"
)

type User struct {
	ID           int64
	Name         string
	CPF          string
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

// ValidateCPF validates an user's CPF, accepts both pure numbers format (12345678912) and formatted (123.456.789-12)
//
// Can return [ErrInvalidCPF] if the CPF is invalid.
func (u *User) ValidateCPF() error {
	u.CPF = reCPF.ReplaceAllString(u.CPF, "")
	if !isValidCPF(u.CPF) {
		return ErrInvalidCPF
	}
	return nil
}

// ValidateEmail validate an user's e-mail, utilizing the [mail.ParseAddress] function.
//
// Can return [ErrInvalidEmail] if the E-mail is invalid.
func (u *User) ValidateEmail() error {
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return ErrInvalidEmail
	}
	return nil
}

var reCPF = regexp.MustCompile(`\D`)

func allSame(s string) bool {
	for i := range len(s) {
		if s[i] != s[0] {
			return false
		}
	}
	return len(s) > 0
}

func isValidCPF(cpf string) bool {
	if len(cpf) != 11 || allSame(cpf) {
		return false
	}
	for i := 9; i < 11; i++ {
		soma := 0
		for j := 0; j < i; j++ {
			soma += int(cpf[j]-'0') * ((i + 1) - j)
		}
		dv := (soma * 10 % 11) % 10
		if dv != int(cpf[i]-'0') {
			return false
		}
	}
	return true
}

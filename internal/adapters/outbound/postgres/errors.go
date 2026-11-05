package postgres

import (
	"errors"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUniqueViolation = errors.New("unique constraint violation")

func translateError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case "23505": // unique_violation
			switch pgErr.ConstraintName {
			case "users_email_key":
				return domain.ErrEmailAlreadyExists
			case "users_cpf_key":
				return domain.ErrCPFAlreadyExists
			}
			return domain.ErrConflict
		}
	}
	return err
}

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
)

var (
	_ ports.UserRepository = (*UserRepo)(nil)
)

type UserRepo struct {
	db *sql.DB
}

// NewUserRepo creates a new Postgres user repository.
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create implements [ports.UserRepository].
//
// Create inserts an user into the DB.
// The only parameters [domain.User] needs are:
//   - [domain.User.Name]
//   - [domain.User.CPF]
//   - [domain.User.Email]
//   - [domain.User.PasswordHash]
//
// and the same user is seeded with the corresponding [domain.User.ID] and [domain.User.CreatedAt].
//
// Can return specific [domain.ErrEmailAlreadyExists].
//
// Can return specific [domain.ErrCPFAlreadyExists].
//
// Can return other DB errors.
func (u *UserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name,cpf,email,password_hash) VALUES ($1,$2,$3,$4) RETURNING id,created_at`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	if err := u.db.QueryRowContext(ctx, query, user.Name, user.CPF, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt); err != nil {
		translated := translateError(err)
		if !errors.Is(err, translated) {
			return translated
		}
		return fmt.Errorf("postgres: create user: %w", err)
	}
	return nil
}

// Read implements [ports.UserRepository].
//
// Read reads an user from the DB by their ID.
//
// The returning user is populated with everything, except [domain.User.PasswordHash].
//
// Returns [domain.ErrNotFound] if user is not found.
//
// Returns specific domain errors in other cases.
func (u *UserRepo) Read(ctx context.Context, userID int64) (domain.User, error) {
	var user domain.User
	query := `SELECT id,name,cpf,email,created_at,updated_at FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	err := u.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Name, &user.CPF, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrNotFound
	}
	if err != nil {
		translated := translateError(err)
		if !errors.Is(err, translated) {
			return domain.User{}, translated
		}
		return domain.User{}, fmt.Errorf("postgres: getByID user: %w", err)
	}
	return user, nil
}

// ReadByEmail implements [ports.UserRepository].
//
// ReadByEmail reads an user by e-mail.
//
// ReadByEmail is used mostly in the login feature, to create the user token.
//
// The [domain.User] returned has only [domain.User.ID], [domain.User.Email] and [domain.User.PasswordHash] populated.
//
// Returns [domain.ErrNotFound] if the user with the corresponding e-mail isn't found in the DB.
//
// Returns specific domain errors in other cases.
func (u *UserRepo) ReadByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	query := `SELECT id,email,password_hash FROM users WHERE email = $1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	err := u.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrNotFound
	}
	if err != nil {
		translated := translateError(err)
		if !errors.Is(err, translated) {
			return domain.User{}, translated
		}
		return domain.User{}, fmt.Errorf("postgres: getByEmail user: %w", err)
	}
	return user, nil
}

// Update implements [ports.UserRepository].
//
// Update updates an user and sets [domain.User.UpdatedAt] to [time.Now] in the DB.
//
// Editable fields:
//
//   - [domain.User.Name]
//   - [domain.User.CPF]
//   - [domain.User.Email]
//   - [domain.User.PasswordHash]
//
// It uses SQL COALESCE, therefore any empty or null (where applicable) fields are not edited.
//
// Returns [domain.ErrNotFound] if no users are edited (ID did not match any users).
//
// Returns specific domain errors in other cases.
func (u *UserRepo) Update(ctx context.Context, user domain.User) error {
	query := `UPDATE users
		SET
			name=COALESCE(NULLIF($1,''),name),
			cpf=COALESCE(NULLIF($2,''),cpf),
			email=COALESCE(NULLIF($3,''),email),
			password_hash=COALESCE($4,password_hash),
			updated_at=NOW()
		WHERE id=$5;
		`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := u.db.ExecContext(ctx, query, user.Name, user.CPF, user.Email, user.PasswordHash, user.ID)
	if err != nil {
		translated := translateError(err)
		if !errors.Is(err, translated) {
			return translated
		}
		return fmt.Errorf("postgres: update users: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rows affected update users: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Delete implements [ports.UserRepository].
//
// Delete deletes an user by their ID in the DB.
//
// Returns [domain.ErrNotFound] if no users are deleted (ID did not match any users).
//
// Returns specific domain errors in other cases.
func (u *UserRepo) Delete(ctx context.Context, userID int64) error {
	query := `DELETE FROM users WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, PostgresQueryTimeout)
	defer cancel()
	res, err := u.db.ExecContext(ctx, query, userID)
	if err != nil {
		translated := translateError(err)
		if !errors.Is(err, translated) {
			return translated
		}
		return fmt.Errorf("postgres: delete user: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres: rows affected delete users: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

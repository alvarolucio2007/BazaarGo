package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("bazaargo/internal/service")

// UserService implements the use cases related to user management and authentication.
type UserService struct {
	repo   ports.UserRepository
	hasher ports.Hasher
	auth   ports.Authenticator
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo ports.UserRepository, hasher ports.Hasher, auth ports.Authenticator) *UserService {
	return &UserService{repo: repo, hasher: hasher, auth: auth}
}

// CreateUserPayload holds the data required to register a new user.
type CreateUserPayload struct {
	Name     string
	CPF      string
	Email    string
	Password string
}

// Create orchestrates user creation.
// It validates the CPF and E-mail, hashes the password, and persists the user in the database.
//
// Returns [domain.ErrInvalidCPF] or [domain.ErrInvalidEmail] for invalid inputs.
// Returns an error if password hashing or database insertion fails.
func (s *UserService) Create(ctx context.Context, payload CreateUserPayload) (domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserService.CreateUser")
	defer span.End()
	user := domain.User{
		Name:  payload.Name,
		CPF:   payload.CPF,
		Email: payload.Email,
	}

	if err := user.ValidateCPF(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.User{}, err
	}
	if err := user.ValidateEmail(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.User{}, err
	}

	hash, err := s.hasher.Hash(payload.Password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.User{}, fmt.Errorf("service: hash password: %w", err)
	}
	user.PasswordHash = []byte(hash)

	if err := s.repo.Create(ctx, &user); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.User{}, err
	}

	return user, nil
}

// Read retrieves a user from the database by their ID.
//
// Returns [domain.ErrNotFound] if the user does not exist.
func (s *UserService) Read(ctx context.Context, userID int64) (domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserService.ReadUser")
	defer span.End()
	user, err := s.repo.Read(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.User{}, err
	}
	return user, nil
}

// ReadByEmail retrieves a user from the database by their E-mail.
//
// Returns [domain.ErrNotFound] if the user does not exist.
func (s *UserService) ReadByEmail(ctx context.Context, email string) (domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserService.ReadUserByEmail")
	defer span.End()
	user, err := s.repo.ReadByEmail(ctx, email)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.User{}, err
	}
	return user, nil
}

// UpdateUserPayload holds the optional fields for updating a user.
// A nil pointer means the field should not be modified.
type UpdateUserPayload struct {
	ID       int64
	Name     *string
	CPF      *string
	Email    *string
	Password *string
}

// Update orchestrates the partial update of a user.
// It validates and hashes data only if the corresponding field is provided.
//
// Returns [domain.ErrInvalidCPF] or [domain.ErrInvalidEmail] if the provided data is invalid.
// Returns an error if password hashing or the database update fails.
func (s *UserService) Update(ctx context.Context, payload UpdateUserPayload) error {
	ctx, span := tracer.Start(ctx, "UserService.UpdateUser")
	defer span.End()
	user := domain.User{ID: payload.ID}

	if payload.Name != nil {
		user.Name = *payload.Name
	}
	if payload.CPF != nil {
		user.CPF = *payload.CPF
		if err := user.ValidateCPF(); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
	}
	if payload.Email != nil {
		user.Email = *payload.Email
		if err := user.ValidateEmail(); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
	}
	if payload.Password != nil {
		hash, err := s.hasher.Hash(*payload.Password)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return fmt.Errorf("service: hash password: %w", err)
		}
		user.PasswordHash = []byte(hash)
	}

	if err := s.repo.Update(ctx, user); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// Delete removes a user from the database by their ID.
//
// Returns [domain.ErrNotFound] if the user does not exist.
func (s *UserService) Delete(ctx context.Context, userID int64) error {
	ctx, span := tracer.Start(ctx, "UserService.DeleteUser")
	defer span.End()
	if err := s.repo.Delete(ctx, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

// VerifyUser handles the authentication flow.
// It verifies the provided credentials and, if valid, generates a PASETO v4 token.
//
// Returns [domain.ErrInvalidCredentials] if the user is not found or the password does not match.
// Returns an error if the token generation fails.
func (s *UserService) VerifyUser(ctx context.Context, email, password string) (string, error) {
	ctx, span := tracer.Start(ctx, "UserService.VerifyUser")
	defer span.End()
	user, err := s.repo.ReadByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			span.SetStatus(codes.Error, "not found")
			return "", domain.ErrInvalidCredentials
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}
	isSame, err := s.hasher.Check(password, string(user.PasswordHash))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", fmt.Errorf("service: check password: %w", err)
	}
	if !isSame {
		span.SetStatus(codes.Error, "invalid credentials")
		return "", domain.ErrInvalidCredentials
	}
	claims, err := ports.NewClaims(email, user.ID, 24*time.Hour)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", fmt.Errorf("service: create claims: %w", err)
	}
	token, err := s.auth.CreateToken(claims)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", fmt.Errorf("service: create token: %w", err)
	}
	return token, nil
}

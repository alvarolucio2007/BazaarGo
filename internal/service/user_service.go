package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
)

type UserService struct {
	repo   ports.UserRepository
	hasher ports.Hasher
	auth   ports.Authenticator
}

func NewUserService(repo ports.UserRepository, hasher ports.Hasher, auth ports.Authenticator) *UserService {
	return &UserService{repo: repo, hasher: hasher, auth: auth}
}

type CreateUserPayload struct {
	Name     string
	CPF      string
	Email    string
	Password string
}
type UserResponse struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CPF       string     `json:"cpf"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// Create orchestrates user creation, by validating the CPF,hashing the password,
// and creating the user in the database.
//
// Returns an error if it detects an invalid CPF.
//
// Returns an error if it detects an invalid E-mail.
//
// Returns an error if the password hashing doesn't complete sucessfully.
//
// Returns an error if the user insertion in the database doesn't complete sucessfully.
func (s *UserService) Create(ctx context.Context, payload CreateUserPayload) (domain.User, error) {
	user := domain.User{
		Name:  payload.Name,
		CPF:   payload.CPF,
		Email: payload.Email,
	}
	if err := user.ValidateCPF(); err != nil {
		return domain.User{}, err
	}
	if err := user.ValidateEmail(); err != nil {
		return domain.User{}, err
	}
	hash, err := s.hasher.Hash(payload.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("service: hash password: %w", err)
	}
	user.PasswordHash = []byte(hash)

	if err := s.repo.Create(ctx, &user); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// Read simply reads an user from the database by ID.
func (s *UserService) Read(ctx context.Context, userID int64) (domain.User, error) {
	return s.repo.Read(ctx, userID)
}

// ReadByEmail simply reads an user from the database by E-mail.
func (s *UserService) ReadByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.repo.ReadByEmail(ctx, email)
}

type UpdateUserPayload struct {
	ID       int64
	Name     *string
	CPF      *string
	Email    *string
	Password *string
}

// Update orchestrates user updating, by validating the CPF,hashing the password,
// and updating the user in the database.
//
// Returns an error if it detects an invalid CPF.
//
// Returns an error if the password hashing doesn't complete sucessfully.
//
// Returns an error if the user update in the database doesn't complete sucessfully.
func (s *UserService) Update(ctx context.Context, payload UpdateUserPayload) error {
	user := domain.User{ID: payload.ID}
	if payload.Name != nil {
		user.Name = *payload.Name
	}
	if payload.CPF != nil {
		user.CPF = *payload.CPF
		if err := user.ValidateCPF(); err != nil {
			return err
		}
	}
	if payload.Email != nil {
		user.Email = *payload.Email
		if err := user.ValidateEmail(); err != nil {
			return err
		}
	}
	if payload.Password != nil {
		hash, err := s.hasher.Hash(*payload.Password)
		if err != nil {
			return fmt.Errorf("service: hash password: %w", err)
		}
		user.PasswordHash = []byte(hash)
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}
	return nil
}

// Delete simply deletes an user from the database by ID.
func (s *UserService) Delete(ctx context.Context, userID int64) error {
	return s.repo.Delete(ctx, userID)
}

// VerifyUser handles user login by checking if the e-mail exists in the database,
// and if the passwords match (via hash comparison). If the credentials are valid,
// it creates the token and returns it.
//
// Returns [domain.ErrInvalidCredentials] if either e-mail doesn't exist or passwords don't match.
//
// Returns respective adapter errors in other cases.
func (s *UserService) VerifyUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.ReadByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}
	isSame, err := s.hasher.Check(password, string(user.PasswordHash))
	if err != nil {
		return "", err
	}
	if !isSame {
		return "", domain.ErrInvalidCredentials
	}
	claims, err := ports.NewClaims(email, user.ID, 24*time.Hour)
	if err != nil {
		return "", err
	}
	token, err := s.auth.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

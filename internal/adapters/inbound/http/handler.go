package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	"github.com/alvarolucio2007/bazaargo/internal/service"
)

var _ StrictServerInterface = (*ServerHandler)(nil)

// ServerHandler is the primary inbound adapter for HTTP requests.
// It implements the generated [StrictServerInterface] and maps HTTP requests
// to the application's core services.
type ServerHandler struct {
	userService *service.UserService
	logger      *slog.Logger
}

// NewServerHandler creates a new instance of ServerHandler with its required dependencies.
func NewServerHandler(userService *service.UserService, logger *slog.Logger) *ServerHandler {
	return &ServerHandler{userService: userService, logger: logger}
}

// CreateToken implements [StrictServerInterface].
//
// It handles the authentication flow by verifying user credentials and
// returning a valid PASETO v4 token.
//
// Returns a 401 Unauthorized response if the e-mail or password are invalid.
//
// Returns a generic internal error (500) on unexpected service failures.
func (s *ServerHandler) CreateToken(ctx context.Context, request CreateTokenRequestObject) (CreateTokenResponseObject, error) {
	token, err := s.userService.VerifyUser(ctx, string(request.Body.Email), request.Body.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return CreateToken401JSONResponse{Error: "invalid email or password"}, nil
		}

		s.logger.Error("unexpected error creating token", "error", err)
		return nil, err
	}

	return CreateToken200JSONResponse{Token: token}, nil
}

// RegisterUser implements [StrictServerInterface].
//
// It handles the user registration flow by passing the request payload
// to the [service.UserService.Create].
//
// Returns a 400 Bad Request response if the CPF or e-mail format is invalid.
//
// Returns a 409 Conflict response if the e-mail or CPF are already in use.
//
// Returns a generic internal error (500) on unexpected service failures.
func (s *ServerHandler) RegisterUser(ctx context.Context, request RegisterUserRequestObject) (RegisterUserResponseObject, error) {
	payload := service.CreateUserPayload{
		Name:     request.Body.Name,
		CPF:      request.Body.Cpf,
		Email:    string(request.Body.Email),
		Password: request.Body.Password,
	}

	user, err := s.userService.Create(ctx, payload)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCPF) {
			return RegisterUser400JSONResponse{Error: "invalid CPF format"}, nil
		}

		if errors.Is(err, domain.ErrInvalidEmail) {
			return RegisterUser400JSONResponse{Error: "invalid Email format"}, nil
		}

		if errors.Is(err, domain.ErrEmailAlreadyExists) || errors.Is(err, domain.ErrCPFAlreadyExists) {
			return RegisterUser409JSONResponse{Error: "user with this email or CPF already exists"}, nil
		}

		s.logger.Error("unexpected error creating user", "error", err)
		return nil, err
	}

	return RegisterUser201JSONResponse{
		Cpf:   user.CPF,
		Email: user.Email,
		Id:    user.ID,
		Name:  user.Name,
	}, nil
}

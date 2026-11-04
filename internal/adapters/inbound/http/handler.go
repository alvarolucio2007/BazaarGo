package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	"github.com/alvarolucio2007/bazaargo/internal/service"
)

var _ StrictServerInterface = (*ServerHandler)(nil)

type ServerHandler struct {
	userService *service.UserService
	logger      *slog.Logger
}

// CreateToken implements [StrictServerInterface].
func (s *ServerHandler) CreateToken(ctx context.Context, request CreateTokenRequestObject) (CreateTokenResponseObject, error) {
	token, err := s.userService.VerifyUser(ctx, string(request.Body.Email), request.Body.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return CreateToken401JSONResponse{Error: "invalid email or password"}, nil
		}
		s.logger.Error("create token error", "error", err)
		return CreateToken500JSONResponse{Error: "internal server error"}, nil
	}

	return CreateToken200JSONResponse{Token: token}, nil
}

// RegisterUser implements [StrictServerInterface].
func (s *ServerHandler) RegisterUser(ctx context.Context, request RegisterUserRequestObject) (RegisterUserResponseObject, error) {
	payload := service.CreateUserPayload{Name: request.Body.Name, CPF: request.Body.Cpf, Email: string(request.Body.Email), Password: request.Body.Password}
	user, err := s.userService.Create(ctx, payload)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCPF) {
			return RegisterUser400JSONResponse{Error: "invalid CPF"}, nil
		}
		if errors.Is(err, domain.ErrInvalidEmail) {
			return RegisterUser400JSONResponse{Error: "invalid Email"}, nil
		}
		if errors.Is(err, domain.ErrEmailAlreadyExists) || errors.Is(err, domain.ErrCPFAlreadyExists) {
			return RegisterUser409JSONResponse{Error: "user already exists"}, nil
		}
		s.logger.Error("unexpected error creating user", "error", err)
		return nil, err
	}
	return RegisterUser201JSONResponse{Cpf: user.CPF, Email: user.Email, Id: user.ID, Name: user.Name}, nil
}

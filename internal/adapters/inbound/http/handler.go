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
	userService    *service.UserService
	productService *service.ProductService
	logger         *slog.Logger
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

// CreateProduct implements [StrictServerInterface].
func (s *ServerHandler) CreateProduct(ctx context.Context, request CreateProductRequestObject) (CreateProductResponseObject, error) {
	payload := service.CreateProductPayload{
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Price:       float64(request.Body.Price),
		Stock:       request.Body.Stock,
		CategoryID:  request.Body.CategoryId,
		Specs:       request.Body.Specs,
	}
	product, err := s.productService.Create(ctx, payload)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidProduct) {
			return CreateProduct400JSONResponse{
				Error: "Invalid product. Please check the credentials.",
			}, nil
		}
		s.logger.ErrorContext(ctx, "unexpected error creating product", "error", err)
		return nil, err
	}
	return CreateProduct201JSONResponse{
		Name:        product.Name,
		Description: product.Description,
		Price:       float32(product.Price),
		Stock:       product.Stock,
		CategoryId:  product.CategoryID,
		Specs:       product.Specs,
	}, nil
}

// GetProductByID implements [StrictServerInterface].
func (s *ServerHandler) GetProductByID(ctx context.Context, request GetProductByIDRequestObject) (GetProductByIDResponseObject, error) {
	id := request.Id
	product, err := s.productService.ReadByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return GetProductByID404JSONResponse{
				Error: "Product not found",
			}, nil
		}
		s.logger.ErrorContext(ctx, "unexpected error creating product", "error", err)
		return nil, err
	}
	return GetProductByID200JSONResponse{
		Name:        product.Name,
		Description: product.Description,
		Price:       float32(product.Price),
		Stock:       product.Stock,
		CategoryId:  product.CategoryID,
		Specs:       product.Specs,
	}, nil
}

// UpdateProduct implements [StrictServerInterface].
func (s *ServerHandler) UpdateProduct(ctx context.Context, request UpdateProductRequestObject) (UpdateProductResponseObject, error) {
	payload := service.UpdateProductPayload{
		ID:          request.Id,
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Price:       float64(request.Body.Price),
		Stock:       request.Body.Stock,
		CategoryID:  request.Body.CategoryId,
		Specs:       request.Body.Specs,
	}
	if err := s.productService.Update(ctx, payload); err != nil {
		if errors.Is(err, domain.ErrInvalidProduct) {
			return UpdateProduct400JSONResponse{
				Error: "Invalid product. Please check the credentials.",
			}, nil
		}
		if errors.Is(err, domain.ErrNotFound) {
			return UpdateProduct404JSONResponse{
				Error: "Not Found",
			}, nil
		}
		s.logger.ErrorContext(ctx, "unexpected error updating product", "error", err)
		return nil, err
	}
	return UpdateProduct204Response{}, nil
}

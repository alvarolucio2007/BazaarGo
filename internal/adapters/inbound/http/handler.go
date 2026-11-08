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
func NewServerHandler(userService *service.UserService, productService *service.ProductService, logger *slog.Logger) *ServerHandler {
	return &ServerHandler{userService: userService, productService: productService, logger: logger}
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
		return CreateToken500JSONResponse{Error: "internal server error"}, nil
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
// It handles the creation of a new product by mapping the request payload to the
// application's core service.
//
// Returns a 400 Bad Request if the product data fails domain validation.
// Returns a 500 Internal Server Error on unexpected failures.
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
				Error: "Invalid product. Please check the product data.",
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
// It fetches a specific product by its ID using the application's core service.
//
// Returns a 404 Not Found if the product does not exist.
// Returns a 500 Internal Server Error on unexpected failures.
func (s *ServerHandler) GetProductByID(ctx context.Context, request GetProductByIDRequestObject) (GetProductByIDResponseObject, error) {
	product, err := s.productService.ReadByID(ctx, request.Id, request.Params.XClientID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return GetProductByID404JSONResponse{Error: "product not found"}, nil
		}

		s.logger.ErrorContext(ctx, "unexpected error fetching product", "error", err)
		return nil, err
	}

	return GetProductByID200JSONResponse{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       float32(product.Price),
		Stock:       product.Stock,
		CategoryId:  product.CategoryID,
		Specs:       product.Specs,
	}, nil
}

// GetProductsSearchName implements [StrictServerInterface].
// It retrieves a list of products matching the provided name query parameter.
//
// Returns a 200 OK response with the list of products (empty if none match).
// Returns a 500 Internal Server Error on unexpected failures.
func (s *ServerHandler) GetProductsSearchName(ctx context.Context, request GetProductsSearchNameRequestObject) (GetProductsSearchNameResponseObject, error) {
	name := request.Name
	products, err := s.productService.ListByName(ctx, name)
	if err != nil {
		s.logger.ErrorContext(ctx, "unexpected error fetching products by name", "error", err)
		return nil, err
	}
	response := make(ListProductResponse, 0, len(products))
	for _, p := range products {
		response = append(response, ProductResponse{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       float32(p.Price),
			Stock:       p.Stock,
			CategoryId:  p.CategoryID,
			Specs:       p.Specs,
		})
	}
	return GetProductsSearchName200JSONResponse{response}, nil
}

// GetRecentProducts implements [StrictServerInterface].
func (s *ServerHandler) GetRecentProducts(ctx context.Context, request GetRecentProductsRequestObject) (GetRecentProductsResponseObject, error) {
	products, err := s.productService.GetRecent(ctx, request.Params.XClientID)
	if err != nil {
		s.logger.ErrorContext(ctx, "unexpected error fetching recent products", "error", err)
		return nil, err
	}

	response := make(ListProductResponse, 0, len(products))
	for _, p := range products {
		response = append(response, ProductResponse{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       float32(p.Price),
			Stock:       p.Stock,
			CategoryId:  p.CategoryID,
			Specs:       p.Specs,
		})
	}

	return GetRecentProducts200JSONResponse(response), nil
}

// UpdateProduct implements [StrictServerInterface].
// It updates an existing product's data based on the provided ID and payload.
//
// Returns a 400 Bad Request if the updated product data fails domain validation.
// Returns a 404 Not Found if the product does not exist.
// Returns a 500 Internal Server Error on unexpected failures.
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
				Error: "Invalid product. Please check the product data.",
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

// DeleteProductsId implements [StrictServerInterface].
// It removes a product from the catalog by its ID.
//
// Returns a 404 Not Found if the product does not exist.
// Returns a 500 Internal Server Error on unexpected failures.
func (s *ServerHandler) DeleteProductsId(ctx context.Context, request DeleteProductsIdRequestObject) (DeleteProductsIdResponseObject, error) {
	id := request.Id

	if err := s.productService.Delete(ctx, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return DeleteProductsId404JSONResponse{
				Error: "Not Found",
			}, nil
		}

		s.logger.ErrorContext(ctx, "unexpected error deleting product", "error", err)
		return nil, err
	}

	return DeleteProductsId204Response{}, nil
}

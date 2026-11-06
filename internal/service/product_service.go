package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/valkey"
	"github.com/alvarolucio2007/bazaargo/internal/domain"
	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
)

// ProductService implements the use cases related to product management.
type ProductService struct {
	repo   ports.ProductRepository
	cache  ports.Cache[domain.Product]
	logger *slog.Logger
}

// NewProductService creates a new instance of ProductService.
func NewProductService(repo ports.ProductRepository, cache ports.Cache[domain.Product], logger *slog.Logger) *ProductService {
	return &ProductService{repo: repo, cache: cache, logger: logger}
}

// CreateProductPayload holds the necessary data to create a new product.
type CreateProductPayload struct {
	Name        string
	Description string
	Price       float64
	Stock       int
	CategoryID  string
	Specs       map[string]any
}

// Create orchestrates the creation of a new product.
// It first persists the product in the database and then attempts to store it in the cache.
//
// Returns the created [domain.Product] populated with the generated ID.
//
// Returns [domain.ErrInvalidProduct] if either name is empty, price is <=0, stock is <0 or category id is empty.
//
// Cache insertion failures are logged but do not prevent the product creation.
func (s *ProductService) Create(ctx context.Context, payload CreateProductPayload) (domain.Product, error) {
	product := domain.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Stock:       payload.Stock,
		CategoryID:  payload.CategoryID,
		Specs:       payload.Specs,
	}
	if err := product.Validate(); err != nil {
		return domain.Product{}, err
	}
	if err := s.repo.Create(ctx, &product); err != nil {
		return domain.Product{}, err
	}

	if err := s.cache.Create(ctx, product); err != nil {
		s.logger.ErrorContext(ctx, "failed to cache product", "error", err)
	}

	return product, nil
}

// ReadByID retrieves a product by its ID using the Cache-Aside pattern.
// It tries to fetch the product from the cache first. If it's a cache miss,
// it fetches from the database and populates the cache for future requests.
//
// Returns [domain.ErrNotFound] if the product doesn't exist in the database.
func (s *ProductService) ReadByID(ctx context.Context, productID string) (domain.Product, error) {
	product, err := s.cache.Read(ctx, productID)
	if err == nil {
		return product, nil
	}
	if !errors.Is(err, valkey.ErrCacheMiss) {
		s.logger.ErrorContext(ctx, "failed to read from cache", "error", err)
	}
	product, err = s.repo.ReadByID(ctx, productID)
	if err != nil {
		return domain.Product{}, err
	}
	_ = s.cache.Create(ctx, product)
	return product, nil
}

// ListByName retrieves a list of products matching the given name.
// It delegates the search directly to the repository.
func (s *ProductService) ListByName(ctx context.Context, name string) ([]domain.Product, error) {
	return s.repo.ListByName(ctx, name)
}

// UpdateProductPayload holds the optional fields for updating a product.
// A nil pointer indicates that the field should not be modified.
type UpdateProductPayload struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Stock       int
	CategoryID  string
	Specs       map[string]any
}

// Update orchestrates the partial update of a product.
// After a successful update in the database, it invalidates the product cache
// to ensure the next read fetches the most recent data.
//
// Returns an error if the repository update fails.
func (s *ProductService) Update(ctx context.Context, payload UpdateProductPayload) error {
	product := domain.Product{
		ID:          payload.ID,
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Stock:       payload.Stock,
		CategoryID:  payload.CategoryID,
		Specs:       payload.Specs,
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return err
	}

	if err := s.cache.Delete(ctx, product.ID); err != nil {
		s.logger.ErrorContext(ctx, "error while deleting cached product", "error", err)
	}

	return nil
}

// Delete removes a product by its ID.
// It deletes the product from the database and invalidates its cache entry.
//
// Returns an error if the repository deletion fails.
func (s *ProductService) Delete(ctx context.Context, productID string) error {
	if err := s.repo.Delete(ctx, productID); err != nil {
		return err
	}

	if err := s.cache.Delete(ctx, productID); err != nil {
		s.logger.ErrorContext(ctx, "error while deleting cached product", "error", err)
	}

	return nil
}

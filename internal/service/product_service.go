package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
	"go.opentelemetry.io/otel/codes"
)

// ProductService implements the use cases related to product management.
type ProductService struct {
	repo       ports.ProductRepository
	cache      ports.Cache[domain.Product]
	recentRepo ports.RecentlyViewedRepository
	logger     *slog.Logger
}

// NewProductService creates a new instance of ProductService.
func NewProductService(repo ports.ProductRepository, cache ports.Cache[domain.Product], recentRepo ports.RecentlyViewedRepository, logger *slog.Logger) *ProductService {
	return &ProductService{repo: repo, cache: cache, recentRepo: recentRepo, logger: logger}
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
	ctx, span := tracer.Start(ctx, "ProductService.CreateProduct")
	defer span.End()
	product := domain.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Stock:       payload.Stock,
		CategoryID:  payload.CategoryID,
		Specs:       payload.Specs,
	}
	if err := product.Validate(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return domain.Product{}, err
	}
	if err := s.repo.Create(ctx, &product); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.Product{}, err
	}

	if err := s.cache.Create(ctx, product); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.logger.ErrorContext(ctx, "failed to cache product", "error", err)
	}

	return product, nil
}

// ReadByID retrieves a product by its ID using the Cache-Aside pattern.
// It tries to fetch the product from the cache first. If it's a cache miss,
// it fetches from the database and populates the cache for future requests.
//
// Returns [domain.ErrNotFound] if the product doesn't exist in the database.
func (s *ProductService) ReadByID(ctx context.Context, productID, clientID string) (domain.Product, error) {
	ctx, span := tracer.Start(ctx, "ProductService.ReadByID")
	defer span.End()

	var product domain.Product
	var err error

	product, err = s.cache.Read(ctx, productID)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			s.logger.ErrorContext(ctx, "failed to read from cache", "error", err)
		}

		product, err = s.repo.ReadByID(ctx, productID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return domain.Product{}, err
		}

		if err := s.cache.Create(ctx, product); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			s.logger.ErrorContext(ctx, "failed to cache product", "error", err)
		}
	}

	if clientID != "" {
		if err := s.recentRepo.Add(ctx, clientID, productID); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			s.logger.ErrorContext(ctx, "failed to add product to recently viewed", "error", err)
		}
	}

	return product, nil
}

// ListByName retrieves a list of products matching the given name.
// It delegates the search directly to the repository.
func (s *ProductService) ListByName(ctx context.Context, name string) ([]domain.Product, error) {
	ctx, span := tracer.Start(ctx, "ProductService.ListByName")
	defer span.End()
	products, err := s.repo.ListByName(ctx, name)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return products, nil
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
	ctx, span := tracer.Start(ctx, "ProductService.Update")
	defer span.End()
	product := domain.Product{
		ID:          payload.ID,
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		Stock:       payload.Stock,
		CategoryID:  payload.CategoryID,
		Specs:       payload.Specs,
	}
	if err := product.Validate(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if err := s.repo.Update(ctx, product); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if err := s.cache.Delete(ctx, product.ID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.logger.ErrorContext(ctx, "error while deleting cached product", "error", err)
	}

	return nil
}

// Delete removes a product by its ID.
// It deletes the product from the database and invalidates its cache entry.
//
// Returns an error if the repository deletion fails.
func (s *ProductService) Delete(ctx context.Context, productID string) error {
	ctx, span := tracer.Start(ctx, "ProductService.Delete")
	defer span.End()
	if err := s.repo.Delete(ctx, productID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if err := s.cache.Delete(ctx, productID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.logger.ErrorContext(ctx, "error while deleting cached product", "error", err)
	}

	return nil
}

// readProduct is a helper that implements the Cache-Aside pattern.
// It attempts to read a product from the cache first. On a cache miss,
// it fetches from the database and populates the cache for future requests.
func (s *ProductService) readProduct(ctx context.Context, productID string) (domain.Product, error) {
	ctx, span := tracer.Start(ctx, "ProductService.readProduct")
	defer span.End()
	product, err := s.cache.Read(ctx, productID)
	if err == nil {
		return product, nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.logger.ErrorContext(ctx, "failed to read from cache", "error", err)
	}
	product, err = s.repo.ReadByID(ctx, productID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.Product{}, err
	}
	if err := s.cache.Create(ctx, product); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		s.logger.ErrorContext(ctx, "failed to cache product", "error", err)
	}
	return product, nil
}

// GetRecent retrieves the list of recently viewed products for a given client.
// It fetches the IDs from the Redis repository and then populates the full
// product data using the cache-aside helper.
//
// Products that no longer exist in the database are silently skipped.
func (s *ProductService) GetRecent(ctx context.Context, clientID string) ([]domain.Product, error) {
	ctx, span := tracer.Start(ctx, "ProductService.GetRecent")
	defer span.End()
	productIDs, err := s.recentRepo.Get(ctx, clientID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	products := make([]domain.Product, 0, len(productIDs))
	for _, id := range productIDs {
		prod, err := s.readProduct(ctx, id)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				continue
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		products = append(products, prod)
	}
	return products, nil
}

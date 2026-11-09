package mongo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/codes"
)

var _ ports.ProductRepository = (*ProductRepo)(nil)

// ProductRepo implements the outbound port for MongoDB product operations.
type ProductRepo struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

// NewProductRepo creates a new instance of ProductRepo.
func NewProductRepo(collection *mongo.Collection, logger *slog.Logger) *ProductRepo {
	return &ProductRepo{collection, logger}
}

type productModel struct {
	ID          bson.ObjectID  `bson:"_id,omitempty"`
	Name        string         `bson:"name"`
	Description string         `bson:"description"`
	Price       float64        `bson:"price"`
	Stock       int            `bson:"stock"`
	CategoryID  string         `bson:"category_id"`
	Specs       map[string]any `bson:"specs"`
	CreatedAt   time.Time      `bson:"created_at"`
	UpdatedAt   time.Time      `bson:"updated_at"`
}

func toModel(p domain.Product) (*productModel, error) {
	var objID bson.ObjectID
	var err error
	if p.ID != "" {
		objID, err = bson.ObjectIDFromHex(p.ID)
		if err != nil {
			return nil, err
		}
	}
	return &productModel{
		ID:          objID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CategoryID:  p.CategoryID,
		Specs:       p.Specs,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

func toDomain(m *productModel) domain.Product {
	return domain.Product{
		ID:          m.ID.Hex(),
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		Stock:       m.Stock,
		CategoryID:  m.CategoryID,
		Specs:       m.Specs,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// Create implements [ports.ProductRepository].
// It inserts a new product document into the MongoDB collection and
// updates the input domain object with the generated ID and timestamps.
func (p *ProductRepo) Create(ctx context.Context, product *domain.Product) error {
	ctx, span := tracer.Start(ctx, "ProductRepo.Create")
	defer span.End()

	timeNow := time.Now()
	product.CreatedAt = timeNow

	model, err := toModel(*product)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("mongo: convert to model: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, MongoQueryTimeout)
	defer cancel()

	result, err := p.collection.InsertOne(ctx, model)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("mongo: insert product: %w", err)
	}

	product.ID = result.InsertedID.(bson.ObjectID).Hex()
	return nil
}

// ReadByID implements [ports.ProductRepository].
// It fetches a product by its ID.
//
// Returns [domain.ErrNotFound] if the product does not exist.
func (p *ProductRepo) ReadByID(ctx context.Context, productID string) (domain.Product, error) {
	ctx, span := tracer.Start(ctx, "ProductRepo.ReadByID")
	defer span.End()

	var model productModel

	objID, err := bson.ObjectIDFromHex(productID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.Product{}, domain.ErrNotFound
	}

	ctx, cancel := context.WithTimeout(ctx, MongoQueryTimeout)
	defer cancel()

	err = p.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			span.SetStatus(codes.Error, err.Error())
			return domain.Product{}, domain.ErrNotFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.Product{}, fmt.Errorf("mongo: find product: %w", err)
	}

	return toDomain(&model), nil
}

// ListByName implements [ports.ProductRepository].
//
// ListByName retrieves a list of products from the MongoDB collection based on a filter.
func (p *ProductRepo) ListByName(ctx context.Context, name string) ([]domain.Product, error) {
	ctx, span := tracer.Start(ctx, "ProductRepo.ListByName")
	defer span.End()

	filterName := bson.Regex{Pattern: name, Options: "i"}
	mongoFilter := bson.M{"name": filterName}

	ctx, cancel := context.WithTimeout(ctx, MongoQueryTimeout)
	defer cancel()

	cursor, err := p.collection.Find(ctx, mongoFilter)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("mongo: find products: %w", err)
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			span.RecordError(err)
			p.logger.ErrorContext(ctx, "error while closing cursor", "error", err)
		}
	}()

	var models []productModel
	if err := cursor.All(ctx, &models); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("mongo: decode products: %w", err)
	}

	products := make([]domain.Product, 0, len(models))
	for _, model := range models {
		products = append(products, toDomain(&model))
	}
	return products, nil
}

// Update implements [ports.ProductRepository].
// It updates an existing product in the database based on its ID.
//
// Returns [domain.ErrNotFound] if the product ID does not exist.
func (p *ProductRepo) Update(ctx context.Context, product domain.Product) error {
	ctx, span := tracer.Start(ctx, "ProductRepo.Update")
	defer span.End()

	objID, err := bson.ObjectIDFromHex(product.ID)
	if err != nil {
		span.RecordError(err) // <-- Adicionado
		span.SetStatus(codes.Error, err.Error())
		return domain.ErrNotFound
	}

	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"category_id": product.CategoryID,
			"specs":       product.Specs,
			"updated_at":  time.Now(),
		},
	}

	ctx, cancel := context.WithTimeout(ctx, MongoQueryTimeout)
	defer cancel()

	result, err := p.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("mongo: update product: %w", err)
	}

	if result.MatchedCount == 0 {
		span.SetStatus(codes.Error, domain.ErrNotFound.Error())
		return domain.ErrNotFound
	}

	return nil
}

// Delete implements [ports.ProductRepository].
// It removes a product from the database by its ID.
//
// Returns [domain.ErrNotFound] if the product ID does not exist.
func (p *ProductRepo) Delete(ctx context.Context, productID string) error {
	// Nome do span corrigido
	ctx, span := tracer.Start(ctx, "ProductRepo.Delete")
	defer span.End()

	objID, err := bson.ObjectIDFromHex(productID)
	if err != nil {
		span.RecordError(err) // <-- Adicionado
		span.SetStatus(codes.Error, err.Error())
		return domain.ErrNotFound
	}

	ctx, cancel := context.WithTimeout(ctx, MongoQueryTimeout)
	defer cancel()

	result, err := p.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("mongo: delete product: %w", err)
	}

	if result.DeletedCount == 0 {
		span.SetStatus(codes.Error, domain.ErrNotFound.Error())
		return domain.ErrNotFound
	}

	return nil
}

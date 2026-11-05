package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/domain"
	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ ports.ProductRepository = (*ProductRepo)(nil)

// ProductRepo implements the outbound port for MongoDB product operations.
type ProductRepo struct {
	collection *mongo.Collection
}

type productModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Price       float64            `bson:"price"`
	Stock       int                `bson:"stock"`
	CategoryID  string             `bson:"category_id"`
	Specs       map[string]any     `bson:"specs"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func toModel(p domain.Product) (*productModel, error) {
	var objID primitive.ObjectID
	var err error
	if p.ID != "" {
		objID, err = primitive.ObjectIDFromHex(p.ID)
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
//
// updates the input domain object with the generated ID and timestamps.
func (p *ProductRepo) Create(ctx context.Context, product *domain.Product) error {
	timeNow := time.Now()
	product.CreatedAt = timeNow

	model, err := toModel(*product)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, MongoQueryTimeout)
	defer cancel()

	result, err := p.collection.InsertOne(ctx, model)
	if err != nil {
		return err
	}

	product.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// ReadByName implements [ports.ProductRepository].
// It fetches a product by its ID.
//
// Returns [domain.ErrNotFound] if the product does not exist.
func (p *ProductRepo) ReadByID(ctx context.Context, productID string) (domain.Product, error) {
	var model productModel

	ctx, cancel := context.WithTimeout(ctx, MongoQueryTimeout)
	defer cancel()

	err := p.collection.FindOne(ctx, bson.M{"_id": productID}).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Product{}, domain.ErrNotFound
		}
		return domain.Product{}, err
	}

	return toDomain(&model), nil
}

// ListByName implements [ports.ProductRepository].
//
// ListByName retrieves a list of products from the MongoDB collection based on a filter.
func (p *ProductRepo) ListByName(ctx context.Context, name string) ([]domain.Product, error) {
	filterName := primitive.Regex{Pattern: name, Options: "i"}
	mongoFilter := bson.M{"name": filterName}
	ctx, cancel := context.WithTimeout(ctx, MongoQueryTimeout)
	defer cancel()
	cursor, err := p.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var models []productModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
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
	objID, err := primitive.ObjectIDFromHex(product.ID)
	if err != nil {
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
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Delete implements [ports.ProductRepository].
// It removes a product from the database by its ID.
//
// Returns [domain.ErrNotFound] if the product ID does not exist.
func (p *ProductRepo) Delete(ctx context.Context, productID string) error {
	objID, err := primitive.ObjectIDFromHex(productID) // <-- Corrigido: Recebe string e converte aqui
	if err != nil {
		return domain.ErrNotFound
	}

	ctx, cancel := context.WithTimeout(ctx, MongoQueryTimeout)
	defer cancel()

	result, err := p.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

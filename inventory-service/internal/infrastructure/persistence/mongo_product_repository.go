package persistence

import (
	"context"
	"errors"
	"time"

	"inventory-service/internal/domain"
	"inventory-service/internal/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoProductRepository struct {
	db *database.MongoDBConnector
}

func NewMongoProductRepository(db *database.MongoDBConnector) *mongoProductRepository {
	return &mongoProductRepository{db: db}
}

func (r *mongoProductRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	now := time.Now()
	productDTO := &database.ProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err := r.db.ProductCollection().InsertOne(ctx, productDTO)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.New("product already exists")
		}
		return nil, err
	}

	return product, nil
}

func (r *mongoProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	var productDTO database.ProductDTO

	filter := bson.M{"_id": id}
	err := r.db.ProductCollection().FindOne(ctx, filter).Decode(&productDTO)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.Product{
		ID:          productDTO.ID,
		Name:        productDTO.Name,
		Description: productDTO.Description,
		Price:       productDTO.Price,
		Stock:       productDTO.Stock,
		CategoryID:  productDTO.CategoryID,
		CreatedAt:   productDTO.CreatedAt,
		UpdatedAt:   productDTO.UpdatedAt,
	}, nil
}

func (r *mongoProductRepository) Update(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	filter := bson.M{"_id": product.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"stock":       product.Stock,
			"category_id": product.CategoryID,
			"updated_at":  time.Now(),
		},
	}

	result, err := r.db.ProductCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, nil
	}

	return product, nil
}

func (r *mongoProductRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	result, err := r.db.ProductCollection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return nil
	}

	return nil
}

func (r *mongoProductRepository) List(ctx context.Context, categoryID string, page, limit int) ([]*domain.Product, int, error) {
	filter := bson.M{}
	if categoryID != "" {
		filter["category_id"] = categoryID
	}

	skip := int64((page - 1) * limit)

	count, err := r.db.ProductCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(skip).
		SetSort(bson.M{"name": 1})

	cursor, err := r.db.ProductCollection().Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var productDTOs []database.ProductDTO
	if err := cursor.All(ctx, &productDTOs); err != nil {
		return nil, 0, err
	}

	products := make([]*domain.Product, len(productDTOs))
	for i, dto := range productDTOs {
		products[i] = &domain.Product{
			ID:          dto.ID,
			Name:        dto.Name,
			Description: dto.Description,
			Price:       dto.Price,
			Stock:       dto.Stock,
			CategoryID:  dto.CategoryID,
			CreatedAt:   dto.CreatedAt,
			UpdatedAt:   dto.UpdatedAt,
		}
	}

	return products, int(count), nil
}

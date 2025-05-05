package persistence

import (
	"context"
	"errors"
	"time"

	"inventory-service/internal/domain"
	"inventory-service/internal/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoCategoryRepository struct {
	db *database.MongoDB
}

func NewMongoCategoryRepository(db *database.MongoDB) *mongoCategoryRepository {
	return &mongoCategoryRepository{db: db}
}

func (r *mongoCategoryRepository) Create(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	now := time.Now()
	categoryDTO := &database.CategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	_, err := r.db.CategoryCollection().InsertOne(ctx, categoryDTO)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.New("category already exists")
		}
		return nil, err
	}

	return category, nil
}

func (r *mongoCategoryRepository) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	var categoryDTO database.CategoryDTO

	filter := bson.M{"_id": id}
	err := r.db.CategoryCollection().FindOne(ctx, filter).Decode(&categoryDTO)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.Category{
		ID:          categoryDTO.ID,
		Name:        categoryDTO.Name,
		Description: categoryDTO.Description,
		CreatedAt:   categoryDTO.CreatedAt,
		UpdatedAt:   categoryDTO.UpdatedAt,
	}, nil
}

func (r *mongoCategoryRepository) Update(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	filter := bson.M{"_id": category.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        category.Name,
			"description": category.Description,
			"updated_at":  time.Now(),
		},
	}

	result, err := r.db.CategoryCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, nil
	}

	return category, nil
}

func (r *mongoCategoryRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	result, err := r.db.CategoryCollection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return nil
	}

	return nil
}

func (r *mongoCategoryRepository) List(ctx context.Context) ([]*domain.Category, error) {
	cursor, err := r.db.CategoryCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categoryDTOs []database.CategoryDTO
	if err := cursor.All(ctx, &categoryDTOs); err != nil {
		return nil, err
	}

	categories := make([]*domain.Category, len(categoryDTOs))
	for i, dto := range categoryDTOs {
		categories[i] = &domain.Category{
			ID:          dto.ID,
			Name:        dto.Name,
			Description: dto.Description,
			CreatedAt:   dto.CreatedAt,
			UpdatedAt:   dto.UpdatedAt,
		}
	}

	return categories, nil
}

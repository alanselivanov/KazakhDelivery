package application

import (
	"context"
	"errors"

	"inventory-service/internal/domain"
	"inventory-service/internal/infrastructure/persistence"
)

type CategoryUseCase struct {
	repo persistence.CategoryRepository
}

func NewCategoryUseCase(repo persistence.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{repo: repo}
}

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	if category.Name == "" {
		return nil, errors.New("category name is required")
	}

	return uc.repo.Create(ctx, category)
}

func (uc *CategoryUseCase) GetCategory(ctx context.Context, id string) (*domain.Category, error) {
	if id == "" {
		return nil, errors.New("category ID is required")
	}

	return uc.repo.GetByID(ctx, id)
}

func (uc *CategoryUseCase) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	if category.ID == "" {
		return nil, errors.New("category ID is required")
	}
	if category.Name == "" {
		return nil, errors.New("category name is required")
	}

	return uc.repo.Update(ctx, category)
}

func (uc *CategoryUseCase) DeleteCategory(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("category ID is required")
	}

	return uc.repo.Delete(ctx, id)
}

func (uc *CategoryUseCase) ListCategories(ctx context.Context) ([]*domain.Category, error) {
	return uc.repo.List(ctx)
}

package application

import (
	"context"
	"errors"

	"inventory-service/internal/domain"
	"inventory-service/internal/infrastructure/persistence"
)

type ProductUseCase struct {
	repo persistence.ProductRepository
}

func NewProductUseCase(repo persistence.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if product.Name == "" {
		return nil, errors.New("product name is required")
	}
	if product.Price <= 0 {
		return nil, errors.New("product price must be positive")
	}
	if product.Stock < 0 {
		return nil, errors.New("product stock cannot be negative")
	}

	return uc.repo.Create(ctx, product)
}

func (uc *ProductUseCase) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	if id == "" {
		return nil, errors.New("product ID is required")
	}

	return uc.repo.GetByID(ctx, id)
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if product.ID == "" {
		return nil, errors.New("product ID is required")
	}
	if product.Name == "" {
		return nil, errors.New("product name is required")
	}

	return uc.repo.Update(ctx, product)
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("product ID is required")
	}

	return uc.repo.Delete(ctx, id)
}

func (uc *ProductUseCase) ListProducts(ctx context.Context, categoryID string, page, limit int) ([]*domain.Product, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return uc.repo.List(ctx, categoryID, page, limit)
}

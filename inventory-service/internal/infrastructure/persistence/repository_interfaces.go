package persistence

import (
	"context"

	"inventory-service/internal/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) (*domain.Category, error)
	GetByID(ctx context.Context, id string) (*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) (*domain.Category, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.Category, error)
}

type CacheableProductRepository interface {
	ProductRepository
	InvalidateCache(ctx context.Context, productID string) error
	InvalidateListCache(ctx context.Context) error
}

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) (*domain.Product, error)
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, categoryID string, page, limit int) ([]*domain.Product, int, error)
}

package persistence

import (
	"context"
	"fmt"
	"inventory-service/internal/domain"
	"inventory-service/internal/infrastructure/cache"
	"log"
)

type redisProductRepository struct {
	repo        *mongoProductRepository
	redisClient *cache.RedisClient
}

type CachedListResult struct {
	Products []*domain.Product `json:"products"`
	Total    int               `json:"total"`
}

func NewRedisProductRepository(repo *mongoProductRepository, redisClient *cache.RedisClient) CacheableProductRepository {
	return &redisProductRepository{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (r *redisProductRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	// Create product in the database
	result, err := r.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	// Invalidate list cache when a new product is added
	if err := r.InvalidateListCache(ctx); err != nil {
		log.Printf("Failed to invalidate list cache: %v", err)
	}

	return result, nil
}

func (r *redisProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	return r.repo.GetByID(ctx, id)
}

func (r *redisProductRepository) Update(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	// Update product in the database
	result, err := r.repo.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	// Invalidate caches
	if err := r.InvalidateCache(ctx, product.ID); err != nil {
		log.Printf("Failed to invalidate product cache: %v", err)
	}

	if err := r.InvalidateListCache(ctx); err != nil {
		log.Printf("Failed to invalidate list cache: %v", err)
	}

	return result, nil
}

func (r *redisProductRepository) Delete(ctx context.Context, id string) error {
	// Delete product from the database
	err := r.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate caches
	if err := r.InvalidateCache(ctx, id); err != nil {
		log.Printf("Failed to invalidate product cache: %v", err)
	}

	if err := r.InvalidateListCache(ctx); err != nil {
		log.Printf("Failed to invalidate list cache: %v", err)
	}

	return nil
}

func (r *redisProductRepository) List(ctx context.Context, categoryID string, page, limit int) ([]*domain.Product, int, error) {
	// Generate cache key based on parameters
	cacheKey := fmt.Sprintf("products:list:%s:%d:%d", categoryID, page, limit)

	// Try to get from cache
	var cachedResult CachedListResult
	found, err := r.redisClient.Get(ctx, cacheKey, &cachedResult)
	if err != nil {
		log.Printf("Cache error: %v", err)
	}

	if found {
		log.Println("Data fetched from cache")
		return cachedResult.Products, cachedResult.Total, nil
	}

	// If not in cache, get from database
	products, total, err := r.repo.List(ctx, categoryID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	log.Println("Data fetched from database")

	// Store in cache for future requests
	cachedResult = CachedListResult{
		Products: products,
		Total:    total,
	}

	if err := r.redisClient.Set(ctx, cacheKey, cachedResult); err != nil {
		log.Printf("Failed to cache product list: %v", err)
	}

	return products, total, nil
}

func (r *redisProductRepository) InvalidateCache(ctx context.Context, productID string) error {
	cacheKey := fmt.Sprintf("products:%s", productID)
	return r.redisClient.Delete(ctx, cacheKey)
}

func (r *redisProductRepository) InvalidateListCache(ctx context.Context) error {
	return r.redisClient.DeleteByPattern(ctx, "products:list:*")
}

package inventoryservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Product represents a product in the inventory
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"category_id"`
}

// CachedListResult represents the cached list of products
type CachedListResult struct {
	Products []*Product `json:"products"`
	Total    int        `json:"total"`
}

// RedisClient interface defines the redis cache operations
type RedisClient interface {
	Get(ctx context.Context, key string, value interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	DeleteByPattern(ctx context.Context, pattern string) error
}

// ProductRepository defines the repository operations for products
type ProductRepository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product *Product) (*Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, categoryID string, page, limit int) ([]*Product, int, error)
}

// CacheableProductRepository is a repository with cache invalidation
type CacheableProductRepository interface {
	ProductRepository
	InvalidateCache(ctx context.Context, productID string) error
	InvalidateListCache(ctx context.Context) error
}

// Mock Redis client
type mockRedisClient struct {
	mock.Mock
}

func (m *mockRedisClient) Get(ctx context.Context, key string, value interface{}) (bool, error) {
	args := m.Called(ctx, key, value)
	return args.Bool(0), args.Error(1)
}

func (m *mockRedisClient) Set(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *mockRedisClient) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *mockRedisClient) DeleteByPattern(ctx context.Context, pattern string) error {
	args := m.Called(ctx, pattern)
	return args.Error(0)
}

// Mock Mongo repository
type mockMongoRepository struct {
	mock.Mock
}

func (m *mockMongoRepository) Create(ctx context.Context, product *Product) (*Product, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Product), args.Error(1)
}

func (m *mockMongoRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Product), args.Error(1)
}

func (m *mockMongoRepository) Update(ctx context.Context, product *Product) (*Product, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Product), args.Error(1)
}

func (m *mockMongoRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockMongoRepository) List(ctx context.Context, categoryID string, page, limit int) ([]*Product, int, error) {
	args := m.Called(ctx, categoryID, page, limit)
	return args.Get(0).([]*Product), args.Int(1), args.Error(2)
}

// RedisProductRepository is a caching layer for products
type RedisProductRepository struct {
	repo        *mockMongoRepository
	redisClient RedisClient
}

func NewRedisProductRepository(repo *mockMongoRepository, redisClient RedisClient) *RedisProductRepository {
	return &RedisProductRepository{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (r *RedisProductRepository) Create(ctx context.Context, product *Product) (*Product, error) {
	result, err := r.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	if err := r.InvalidateListCache(ctx); err != nil {
	}

	return result, nil
}

func (r *RedisProductRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	return r.repo.GetByID(ctx, id)
}

func (r *RedisProductRepository) Update(ctx context.Context, product *Product) (*Product, error) {
	result, err := r.repo.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	if err := r.InvalidateCache(ctx, product.ID); err != nil {
	}

	if err := r.InvalidateListCache(ctx); err != nil {
	}

	return result, nil
}

func (r *RedisProductRepository) Delete(ctx context.Context, id string) error {
	err := r.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	if err := r.InvalidateCache(ctx, id); err != nil {
	}

	if err := r.InvalidateListCache(ctx); err != nil {
	}

	return nil
}

func (r *RedisProductRepository) List(ctx context.Context, categoryID string, page, limit int) ([]*Product, int, error) {
	cacheKey := fmt.Sprintf("products:list:%s:%d:%d", categoryID, page, limit)

	var cachedResult CachedListResult
	found, err := r.redisClient.Get(ctx, cacheKey, &cachedResult)
	if err == nil && found {
		return cachedResult.Products, cachedResult.Total, nil
	}

	products, total, err := r.repo.List(ctx, categoryID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	cachedResult = CachedListResult{
		Products: products,
		Total:    total,
	}

	_ = r.redisClient.Set(ctx, cacheKey, cachedResult)

	return products, total, nil
}

func (r *RedisProductRepository) InvalidateCache(ctx context.Context, productID string) error {
	cacheKey := "products:" + productID
	return r.redisClient.Delete(ctx, cacheKey)
}

func (r *RedisProductRepository) InvalidateListCache(ctx context.Context) error {
	return r.redisClient.DeleteByPattern(ctx, "products:list:*")
}

func TestRedisProductRepository_List_CacheHit(t *testing.T) {
	ctx := context.Background()
	mockMongo := new(mockMongoRepository)
	mockRedis := new(mockRedisClient)

	expectedProducts := []*Product{
		{
			ID:          "1",
			Name:        "Test Product 1",
			Description: "Description 1",
			Price:       99.99,
			CategoryID:  "category1",
		},
	}
	expectedTotal := 1

	cacheKey := "products:list:category1:1:10"
	mockRedis.On("Get", ctx, cacheKey, mock.AnythingOfType("*inventoryservice_test.CachedListResult")).
		Run(func(args mock.Arguments) {
			result := args.Get(2).(*CachedListResult)
			result.Products = expectedProducts
			result.Total = expectedTotal
		}).
		Return(true, nil)

	repo := NewRedisProductRepository(mockMongo, mockRedis)

	products, total, err := repo.List(ctx, "category1", 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, expectedTotal, total)
	assert.Equal(t, expectedProducts, products)
	mockRedis.AssertExpectations(t)
	mockMongo.AssertNotCalled(t, "List")
}

func TestRedisProductRepository_List_CacheMiss(t *testing.T) {
	ctx := context.Background()
	mockMongo := new(mockMongoRepository)
	mockRedis := new(mockRedisClient)

	expectedProducts := []*Product{
		{
			ID:          "1",
			Name:        "Test Product 1",
			Description: "Description 1",
			Price:       99.99,
			CategoryID:  "category1",
		},
	}
	expectedTotal := 1

	cacheKey := "products:list:category1:1:10"
	mockRedis.On("Get", ctx, cacheKey, mock.AnythingOfType("*inventoryservice_test.CachedListResult")).Return(false, nil)

	mockMongo.On("List", ctx, "category1", 1, 10).Return(expectedProducts, expectedTotal, nil)

	mockRedis.On("Set", ctx, cacheKey, mock.Anything).Return(nil)

	repo := NewRedisProductRepository(mockMongo, mockRedis)

	products, total, err := repo.List(ctx, "category1", 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, expectedTotal, total)
	assert.Equal(t, expectedProducts, products)
	mockRedis.AssertExpectations(t)
	mockMongo.AssertExpectations(t)
}

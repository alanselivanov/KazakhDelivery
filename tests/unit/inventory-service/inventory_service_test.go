package inventoryservice_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type InventoryItem struct {
	ProductID  string    `json:"product_id"`
	Stock      int       `json:"stock"`
	Reserved   int       `json:"reserved"`
	Available  int       `json:"available"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastSoldAt time.Time `json:"last_sold_at,omitempty"`
}

type InventoryRepository interface {
	GetItem(ctx context.Context, productID string) (*InventoryItem, error)
	UpdateStock(ctx context.Context, productID string, quantity int) error
	UpdateReservation(ctx context.Context, productID string, quantity int) error
	ListLowStock(ctx context.Context, threshold int) ([]*InventoryItem, error)
	BatchUpdateStock(ctx context.Context, updates map[string]int) error
}

type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) GetItem(ctx context.Context, productID string) (*InventoryItem, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*InventoryItem), args.Error(1)
}

func (m *MockInventoryRepository) UpdateStock(ctx context.Context, productID string, quantity int) error {
	args := m.Called(ctx, productID, quantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) UpdateReservation(ctx context.Context, productID string, quantity int) error {
	args := m.Called(ctx, productID, quantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) ListLowStock(ctx context.Context, threshold int) ([]*InventoryItem, error) {
	args := m.Called(ctx, threshold)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*InventoryItem), args.Error(1)
}

func (m *MockInventoryRepository) BatchUpdateStock(ctx context.Context, updates map[string]int) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

type InventoryService struct {
	repo InventoryRepository
}

func NewInventoryService(repo InventoryRepository) *InventoryService {
	return &InventoryService{
		repo: repo,
	}
}

func (s *InventoryService) CheckStock(ctx context.Context, productID string, quantity int) (bool, error) {
	item, err := s.repo.GetItem(ctx, productID)
	if err != nil {
		return false, err
	}

	return item.Available >= quantity, nil
}

func (s *InventoryService) ReserveStock(ctx context.Context, productID string, quantity int) error {
	item, err := s.repo.GetItem(ctx, productID)
	if err != nil {
		return err
	}

	if item.Available < quantity {
		return errors.New("insufficient stock available")
	}

	return s.repo.UpdateReservation(ctx, productID, quantity)
}

func (s *InventoryService) ReleaseStock(ctx context.Context, productID string, quantity int) error {
	return s.repo.UpdateReservation(ctx, productID, -quantity)
}

func (s *InventoryService) ConsumeStock(ctx context.Context, productID string, quantity int) error {
	item, err := s.repo.GetItem(ctx, productID)
	if err != nil {
		return err
	}

	if item.Stock < quantity {
		return errors.New("insufficient stock available")
	}

	return s.repo.UpdateStock(ctx, productID, -quantity)
}

func (s *InventoryService) GetLowStockItems(ctx context.Context, threshold int) ([]*InventoryItem, error) {
	return s.repo.ListLowStock(ctx, threshold)
}

func (s *InventoryService) BatchUpdate(ctx context.Context, updates map[string]int) error {
	return s.repo.BatchUpdateStock(ctx, updates)
}

func TestInventoryService_CheckStock(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockInventoryRepository)
	service := NewInventoryService(mockRepo)

	t.Run("Sufficient Stock", func(t *testing.T) {
		item := &InventoryItem{
			ProductID: "prod-1",
			Stock:     100,
			Reserved:  20,
			Available: 80,
		}

		mockRepo.On("GetItem", ctx, "prod-1").Return(item, nil).Once()

		available, err := service.CheckStock(ctx, "prod-1", 50)

		assert.NoError(t, err)
		assert.True(t, available)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Insufficient Stock", func(t *testing.T) {
		item := &InventoryItem{
			ProductID: "prod-2",
			Stock:     30,
			Reserved:  20,
			Available: 10,
		}

		mockRepo.On("GetItem", ctx, "prod-2").Return(item, nil).Once()

		available, err := service.CheckStock(ctx, "prod-2", 20)

		assert.NoError(t, err)
		assert.False(t, available)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Item Not Found", func(t *testing.T) {
		mockRepo.On("GetItem", ctx, "non-existent").Return(nil, errors.New("item not found")).Once()

		available, err := service.CheckStock(ctx, "non-existent", 10)

		assert.Error(t, err)
		assert.False(t, available)
		mockRepo.AssertExpectations(t)
	})
}

func TestInventoryService_ReserveStock(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockInventoryRepository)
	service := NewInventoryService(mockRepo)

	t.Run("Successful Reservation", func(t *testing.T) {
		item := &InventoryItem{
			ProductID: "prod-1",
			Stock:     100,
			Reserved:  20,
			Available: 80,
		}

		mockRepo.On("GetItem", ctx, "prod-1").Return(item, nil).Once()
		mockRepo.On("UpdateReservation", ctx, "prod-1", 30).Return(nil).Once()

		err := service.ReserveStock(ctx, "prod-1", 30)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Insufficient Stock for Reservation", func(t *testing.T) {
		item := &InventoryItem{
			ProductID: "prod-2",
			Stock:     50,
			Reserved:  40,
			Available: 10,
		}

		mockRepo.On("GetItem", ctx, "prod-2").Return(item, nil).Once()

		err := service.ReserveStock(ctx, "prod-2", 20)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient stock")
		mockRepo.AssertNotCalled(t, "UpdateReservation")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		item := &InventoryItem{
			ProductID: "prod-3",
			Stock:     100,
			Reserved:  0,
			Available: 100,
		}

		mockRepo.On("GetItem", ctx, "prod-3").Return(item, nil).Once()
		mockRepo.On("UpdateReservation", ctx, "prod-3", 50).
			Return(errors.New("database error")).Once()

		err := service.ReserveStock(ctx, "prod-3", 50)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestInventoryService_ReleaseStock(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockInventoryRepository)
	service := NewInventoryService(mockRepo)

	t.Run("Successful Release", func(t *testing.T) {
		mockRepo.On("UpdateReservation", ctx, "prod-1", -20).Return(nil).Once()

		err := service.ReleaseStock(ctx, "prod-1", 20)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.On("UpdateReservation", ctx, "prod-error", -10).
			Return(errors.New("database error")).Once()

		err := service.ReleaseStock(ctx, "prod-error", 10)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestInventoryService_ConsumeStock(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockInventoryRepository)
	service := NewInventoryService(mockRepo)

	t.Run("Successful Consumption", func(t *testing.T) {
		item := &InventoryItem{
			ProductID: "prod-1",
			Stock:     50,
			Reserved:  30,
			Available: 20,
		}

		mockRepo.On("GetItem", ctx, "prod-1").Return(item, nil).Once()
		mockRepo.On("UpdateStock", ctx, "prod-1", -20).Return(nil).Once()

		err := service.ConsumeStock(ctx, "prod-1", 20)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Insufficient Stock", func(t *testing.T) {
		item := &InventoryItem{
			ProductID: "prod-2",
			Stock:     10,
			Reserved:  5,
			Available: 5,
		}

		mockRepo.On("GetItem", ctx, "prod-2").Return(item, nil).Once()

		err := service.ConsumeStock(ctx, "prod-2", 15)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient stock")
		mockRepo.AssertNotCalled(t, "UpdateStock")
		mockRepo.AssertExpectations(t)
	})
}

func TestInventoryService_GetLowStockItems(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockInventoryRepository)
	service := NewInventoryService(mockRepo)

	lowStockItems := []*InventoryItem{
		{
			ProductID: "prod-low-1",
			Stock:     5,
			Reserved:  2,
			Available: 3,
		},
		{
			ProductID: "prod-low-2",
			Stock:     8,
			Reserved:  3,
			Available: 5,
		},
	}

	t.Run("Successfully Get Low Stock Items", func(t *testing.T) {
		mockRepo.On("ListLowStock", ctx, 10).Return(lowStockItems, nil).Once()

		items, err := service.GetLowStockItems(ctx, 10)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(items))
		assert.Equal(t, "prod-low-1", items[0].ProductID)
		assert.Equal(t, "prod-low-2", items[1].ProductID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("No Low Stock Items", func(t *testing.T) {
		mockRepo.On("ListLowStock", ctx, 5).Return([]*InventoryItem{}, nil).Once()

		items, err := service.GetLowStockItems(ctx, 5)

		assert.NoError(t, err)
		assert.Empty(t, items)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.On("ListLowStock", ctx, 20).Return(nil, errors.New("database error")).Once()

		items, err := service.GetLowStockItems(ctx, 20)

		assert.Error(t, err)
		assert.Nil(t, items)
		mockRepo.AssertExpectations(t)
	})
}

func TestInventoryService_BatchUpdate(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockInventoryRepository)
	service := NewInventoryService(mockRepo)

	t.Run("Successful Batch Update", func(t *testing.T) {
		updates := map[string]int{
			"prod-1": 50,
			"prod-2": 20,
			"prod-3": -10,
		}

		mockRepo.On("BatchUpdateStock", ctx, updates).Return(nil).Once()

		err := service.BatchUpdate(ctx, updates)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		updates := map[string]int{
			"prod-1": 50,
			"prod-2": 20,
		}

		mockRepo.On("BatchUpdateStock", ctx, updates).
			Return(errors.New("transaction failed")).Once()

		err := service.BatchUpdate(ctx, updates)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction failed")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty Updates Map", func(t *testing.T) {
		updates := map[string]int{}

		mockRepo.On("BatchUpdateStock", ctx, updates).Return(nil).Once()

		err := service.BatchUpdate(ctx, updates)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

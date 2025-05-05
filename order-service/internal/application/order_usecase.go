package application

import (
	"context"
	"errors"

	"order-service/internal/domain"
	"order-service/internal/infrastructure/persistence"
)

type OrderUseCase struct {
	repo persistence.OrderRepository
}

func NewOrderUseCase(repo persistence.OrderRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	if order.UserID == "" {
		return nil, errors.New("user ID is required")
	}
	if len(order.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	for _, item := range order.Items {
		if item.ProductID == "" {
			return nil, errors.New("product ID is required for all items")
		}
		if item.Quantity <= 0 {
			return nil, errors.New("item quantity must be positive")
		}
		if item.Price <= 0 {
			return nil, errors.New("item price must be positive")
		}
	}

	return uc.repo.Create(ctx, order)
}

func (uc *OrderUseCase) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}

	return uc.repo.GetByID(ctx, id)
}

func (uc *OrderUseCase) UpdateOrderStatus(ctx context.Context, id string, status domain.OrderStatus) (*domain.Order, error) {
	if id == "" {
		return nil, errors.New("order ID is required")
	}

	validStatuses := map[domain.OrderStatus]bool{
		domain.OrderStatusPending:   true,
		domain.OrderStatusCompleted: true,
		domain.OrderStatusCancelled: true,
	}

	if !validStatuses[status] {
		return nil, errors.New("invalid order status")
	}

	order, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, nil
	}

	order.UpdateStatus(status)

	return uc.repo.Update(ctx, order)
}

func (uc *OrderUseCase) ListUserOrders(ctx context.Context, userID string) ([]*domain.Order, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return uc.repo.ListByUserID(ctx, userID)
}

func (uc *OrderUseCase) CheckStock(ctx context.Context, productID string, quantity int) (bool, error) {
	return true, nil
}

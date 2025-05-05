package persistence

import (
	"context"
	"order-service/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) (*domain.Order, error)
	ListByUserID(ctx context.Context, userID string) ([]*domain.Order, error)
}

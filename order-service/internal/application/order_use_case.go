package application

import (
	"context"
	"order-service/internal/domain"
	"order-service/internal/infrastructure/messaging"
	"order-service/internal/infrastructure/persistence"
	"time"
)

type OrderUseCase struct {
	orderRepo      persistence.OrderRepository
	eventPublisher messaging.EventPublisher
}

func NewOrderUseCase(orderRepo persistence.OrderRepository, eventPublisher messaging.EventPublisher) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
	}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, userID string, items []domain.OrderItem) (*domain.Order, error) {

	order := domain.NewOrder(userID, items, domain.OrderStatusPending)

	savedOrder, err := uc.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	go uc.publishOrderCreatedEvent(savedOrder)

	return savedOrder, nil
}

func (uc *OrderUseCase) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	return uc.orderRepo.GetByID(ctx, id)
}

func (uc *OrderUseCase) UpdateOrderStatus(ctx context.Context, id string, status domain.OrderStatus) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, nil
	}

	order.UpdateStatus(status)
	return uc.orderRepo.Update(ctx, order)
}

func (uc *OrderUseCase) ListOrdersByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	return uc.orderRepo.ListByUserID(ctx, userID)
}

func (uc *OrderUseCase) publishOrderCreatedEvent(order *domain.Order) {

	items := make([]messaging.OrderItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = messaging.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	event := messaging.OrderCreatedEvent{
		OrderID:   order.ID,
		UserID:    order.UserID,
		Items:     items,
		Total:     order.Total,
		Timestamp: time.Now().UnixNano(),
	}

	err := uc.eventPublisher.PublishOrderCreated(event)
	if err != nil {
	}
}

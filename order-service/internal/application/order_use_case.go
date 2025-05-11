package application

import (
	"context"
	"fmt"
	"log"
	"order-service/internal/domain"
	"order-service/internal/infrastructure/database"
	"order-service/internal/infrastructure/messaging"
	"order-service/internal/infrastructure/persistence"
	"time"

	"github.com/redis/go-redis/v9"
)

type OrderUseCase struct {
	orderRepo      persistence.OrderRepository
	eventPublisher messaging.EventPublisher
	cache          *database.RedisCache
}

func NewOrderUseCase(orderRepo persistence.OrderRepository, eventPublisher messaging.EventPublisher, cache *database.RedisCache) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
		cache:          cache,
	}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, userID string, items []domain.OrderItem) (*domain.Order, error) {

	order := domain.NewOrder(userID, items, domain.OrderStatusPending)

	savedOrder, err := uc.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	if uc.cache != nil {
		cacheKey := fmt.Sprintf("user_orders:%s", userID)
		if err := uc.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Failed to invalidate cache for user %s: %v", userID, err)
		} else {
			log.Printf("Cache invalidated for user %s", userID)
		}
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
	updatedOrder, err := uc.orderRepo.Update(ctx, order)
	if err != nil {
		return nil, err
	}

	if uc.cache != nil {
		cacheKey := fmt.Sprintf("user_orders:%s", order.UserID)
		if err := uc.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Failed to invalidate cache for user %s: %v", order.UserID, err)
		} else {
			log.Printf("Cache invalidated for user %s", order.UserID)
		}
	}

	return updatedOrder, nil
}

func (uc *OrderUseCase) ListOrdersByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	cacheKey := fmt.Sprintf("user_orders:%s", userID)
	var orders []*domain.Order

	if uc.cache != nil {
		err := uc.cache.Get(ctx, cacheKey, &orders)
		if err == nil {
			log.Printf("Data retrieved from cache for user %s", userID)
			return orders, nil
		} else if err != redis.Nil {
			log.Printf("Redis error: %v", err)
		}
	}

	dbOrders, err := uc.orderRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	log.Printf("Data retrieved from database for user %s", userID)

	if uc.cache != nil && len(dbOrders) > 0 {
		if err := uc.cache.Set(ctx, cacheKey, dbOrders); err != nil {
			log.Printf("Failed to cache orders for user %s: %v", userID, err)
		}
	}

	return dbOrders, nil
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

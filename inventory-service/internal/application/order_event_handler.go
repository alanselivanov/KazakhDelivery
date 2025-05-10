package application

import (
	"context"
	"errors"
	"inventory-service/internal/infrastructure/messaging"
	grpc "inventory-service/internal/infrastructure/product"
	"log"
	"sync"
)

type OrderEventHandler struct {
	productClient grpc.ProductServiceClient
	metrics       *Metrics
}

func NewOrderEventHandler(productClient grpc.ProductServiceClient, metrics *Metrics) *OrderEventHandler {
	return &OrderEventHandler{
		productClient: productClient,
		metrics:       metrics,
	}
}

func (h *OrderEventHandler) HandleOrderCreated(ctx context.Context, event *messaging.OrderCreatedEvent) error {
	log.Printf("Processing order.created event for order ID: %s with %d items",
		event.OrderID, len(event.Items))

	var wg sync.WaitGroup
	var mu sync.Mutex
	var failedItems []messaging.StockUpdateResult

	for _, item := range event.Items {
		wg.Add(1)
		go func(item messaging.OrderItem) {
			defer wg.Done()

			err := h.productClient.DecreaseStock(ctx, item.ProductID, item.Quantity)

			result := messaging.StockUpdateResult{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Success:   err == nil,
			}

			if err != nil {
				result.Error = err.Error()
				log.Printf("Failed to update stock for product %s: %v", item.ProductID, err)

				mu.Lock()
				failedItems = append(failedItems, result)
				mu.Unlock()

				h.metrics.IncStockUpdateErrors()
			} else {
				log.Printf("Successfully updated stock for product %s by %d", item.ProductID, item.Quantity)
			}
		}(item)
	}

	wg.Wait()

	if len(failedItems) > 0 {
		log.Printf("%d items failed stock update for order %s", len(failedItems), event.OrderID)
		return errors.New("some items failed stock update")
	}

	log.Printf("Successfully processed order %s", event.OrderID)
	return nil
}

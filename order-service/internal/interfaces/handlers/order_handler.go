package handlers

import (
	"context"
	"log"

	"order-service/internal/application"
	"order-service/internal/domain"
	"proto/order"
)

type OrderHandler struct {
	order.UnimplementedOrderServiceServer
	orderUseCase *application.OrderUseCase
}

func NewOrderHandler(orderUseCase *application.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *order.OrderRequest) (*order.OrderResponse, error) {
	items := make([]domain.OrderItem, len(req.Order.Items))
	for i, item := range req.Order.Items {
		items[i] = domain.OrderItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			Price:     float64(item.Price),
		}
	}

	createdOrder, err := h.orderUseCase.CreateOrder(ctx, req.Order.UserId, items)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		return nil, err
	}

	return &order.OrderResponse{
		Order: &order.Order{
			Id:     createdOrder.ID,
			UserId: createdOrder.UserID,
			Items:  convertToProtoItems(createdOrder.Items),
			Total:  float32(createdOrder.Total),
			Status: string(createdOrder.Status),
		},
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *order.OrderID) (*order.OrderResponse, error) {
	domainOrder, err := h.orderUseCase.GetOrderByID(ctx, req.Id)
	if err != nil {
		log.Printf("Error getting order: %v", err)
		return nil, err
	}

	if domainOrder == nil {
		return &order.OrderResponse{}, nil
	}

	return &order.OrderResponse{
		Order: &order.Order{
			Id:     domainOrder.ID,
			UserId: domainOrder.UserID,
			Items:  convertToProtoItems(domainOrder.Items),
			Total:  float32(domainOrder.Total),
			Status: string(domainOrder.Status),
		},
	}, nil
}

func (h *OrderHandler) UpdateOrder(ctx context.Context, req *order.OrderRequest) (*order.OrderResponse, error) {
	domainOrder, err := h.orderUseCase.UpdateOrderStatus(ctx, req.Order.Id, domain.OrderStatus(req.Order.Status))
	if err != nil {
		log.Printf("Error updating order: %v", err)
		return nil, err
	}

	if domainOrder == nil {
		return &order.OrderResponse{}, nil
	}

	return &order.OrderResponse{
		Order: &order.Order{
			Id:     domainOrder.ID,
			UserId: domainOrder.UserID,
			Items:  convertToProtoItems(domainOrder.Items),
			Total:  float32(domainOrder.Total),
			Status: string(domainOrder.Status),
		},
	}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *order.UserID) (*order.OrderListResponse, error) {
	orders, err := h.orderUseCase.ListOrdersByUserID(ctx, req.Id)
	if err != nil {
		log.Printf("Error listing orders: %v", err)
		return nil, err
	}

	protoOrders := make([]*order.Order, len(orders))
	for i, domainOrder := range orders {
		protoOrders[i] = &order.Order{
			Id:     domainOrder.ID,
			UserId: domainOrder.UserID,
			Items:  convertToProtoItems(domainOrder.Items),
			Total:  float32(domainOrder.Total),
			Status: string(domainOrder.Status),
		}
	}

	return &order.OrderListResponse{
		Orders: protoOrders,
	}, nil
}

func (h *OrderHandler) CheckStock(ctx context.Context, req *order.StockCheckRequest) (*order.StockCheckResponse, error) {
	log.Printf("Checking stock for product ID: %s, quantity: %d", req.ProductId, req.Quantity)

	return &order.StockCheckResponse{
		Available: true,
	}, nil
}

func convertToProtoItems(items []domain.OrderItem) []*order.OrderItem {
	protoItems := make([]*order.OrderItem, len(items))
	for i, item := range items {
		protoItems[i] = &order.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     float32(item.Price),
		}
	}
	return protoItems
}

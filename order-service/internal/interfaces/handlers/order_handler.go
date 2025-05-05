package handlers

import (
	"context"
	"time"

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
	var items []domain.OrderItem
	for _, item := range req.Order.Items {
		items = append(items, domain.OrderItem{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			Price:     float64(item.Price),
		})
	}

	domainOrder := domain.NewOrder(
		req.Order.UserId,
		items,
		domain.OrderStatus(req.Order.Status),
	)

	created, err := h.orderUseCase.CreateOrder(ctx, domainOrder)
	if err != nil {
		return nil, err
	}

	return convertToOrderResponse(created), nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *order.OrderID) (*order.OrderResponse, error) {
	domainOrder, err := h.orderUseCase.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if domainOrder == nil {
		return nil, nil
	}

	return convertToOrderResponse(domainOrder), nil
}

func (h *OrderHandler) UpdateOrder(ctx context.Context, req *order.OrderRequest) (*order.OrderResponse, error) {
	domainOrder, err := h.orderUseCase.GetOrder(ctx, req.Order.Id)
	if err != nil {
		return nil, err
	}
	if domainOrder == nil {
		return nil, nil
	}

	updated, err := h.orderUseCase.UpdateOrderStatus(ctx, req.Order.Id, domain.OrderStatus(req.Order.Status))
	if err != nil {
		return nil, err
	}

	return convertToOrderResponse(updated), nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *order.UserID) (*order.OrderListResponse, error) {
	domainOrders, err := h.orderUseCase.ListUserOrders(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var protoOrders []*order.Order
	for _, o := range domainOrders {
		protoOrders = append(protoOrders, convertToProtoOrder(o))
	}

	return &order.OrderListResponse{
		Orders: protoOrders,
	}, nil
}

func (h *OrderHandler) CheckStock(ctx context.Context, req *order.StockCheckRequest) (*order.StockCheckResponse, error) {
	available, err := h.orderUseCase.CheckStock(ctx, req.ProductId, int(req.Quantity))
	if err != nil {
		return nil, err
	}

	return &order.StockCheckResponse{
		Available: available,
	}, nil
}

func convertToOrderResponse(domainOrder *domain.Order) *order.OrderResponse {
	return &order.OrderResponse{
		Order: convertToProtoOrder(domainOrder),
	}
}

func convertToProtoOrder(domainOrder *domain.Order) *order.Order {
	var items []*order.OrderItem
	for _, item := range domainOrder.Items {
		items = append(items, &order.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     float32(item.Price),
		})
	}

	return &order.Order{
		Id:        domainOrder.ID,
		UserId:    domainOrder.UserID,
		Items:     items,
		Total:     float32(domainOrder.Total),
		Status:    string(domainOrder.Status),
		CreatedAt: domainOrder.CreatedAt.Format(time.RFC3339),
		UpdatedAt: domainOrder.UpdatedAt.Format(time.RFC3339),
	}
}

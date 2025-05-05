package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

type Order struct {
	ID        string
	UserID    string
	Items     []OrderItem
	Total     float64
	Status    OrderStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOrder(userID string, items []OrderItem, status OrderStatus) *Order {
	if status == "" {
		status = OrderStatusPending
	}

	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}

	return &Order{
		ID:        uuid.New().String(),
		UserID:    userID,
		Items:     items,
		Total:     total,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (o *Order) UpdateStatus(newStatus OrderStatus) {
	o.Status = newStatus
	o.UpdatedAt = time.Now()
}

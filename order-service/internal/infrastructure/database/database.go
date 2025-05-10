package database

import (
	"sync"
	"time"
)

type OrderItemDTO struct {
	ProductID string  `bson:"product_id"`
	Quantity  int     `bson:"quantity"`
	Price     float64 `bson:"price"`
}

type OrderDTO struct {
	ID        string         `bson:"_id,omitempty"`
	UserID    string         `bson:"user_id"`
	Items     []OrderItemDTO `bson:"items"`
	Total     float64        `bson:"total"`
	Status    string         `bson:"status"`
	CreatedAt time.Time      `bson:"created_at"`
	UpdatedAt time.Time      `bson:"updated_at"`
}

type InMemoryDB struct {
	Orders map[string]*OrderDTO
	mu     sync.RWMutex
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		Orders: make(map[string]*OrderDTO),
	}
}

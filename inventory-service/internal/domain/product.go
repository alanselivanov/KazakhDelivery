package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Stock       int
	CategoryID  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProduct(name, description string, price float64, stock int, categoryID string) *Product {
	now := time.Now()
	return &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		CategoryID:  categoryID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

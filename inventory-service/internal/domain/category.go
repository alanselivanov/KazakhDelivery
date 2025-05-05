package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewCategory(name, description string) *Category {
	now := time.Now()
	return &Category{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

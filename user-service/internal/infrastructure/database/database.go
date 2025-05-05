package database

import (
	"sync"
	"time"
)

type UserDTO struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type InMemoryDB struct {
	Users map[string]*UserDTO
	mu    sync.RWMutex
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		Users: make(map[string]*UserDTO),
	}
}

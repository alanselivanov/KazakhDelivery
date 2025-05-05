package database

import (
	"sync"
	"time"
)

type UserDTO struct {
	ID        string    `bson:"_id,omitempty"`
	Username  string    `bson:"username"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
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

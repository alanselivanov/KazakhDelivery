package persistence

import (
	"context"
	"sync"

	"user-service/internal/domain"
	"user-service/internal/infrastructure/database"
)

type userRepository struct {
	db *database.InMemoryDB
	mu sync.RWMutex
}

func NewUserRepository(db *database.InMemoryDB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	dto := &database.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}

	r.db.Users[user.ID] = dto

	return &domain.User{
		ID:        dto.ID,
		Username:  dto.Username,
		Email:     dto.Email,
		Password:  dto.Password,
		CreatedAt: dto.CreatedAt,
	}, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	dto, exists := r.db.Users[id]
	if !exists {
		return nil, nil
	}

	return &domain.User{
		ID:        dto.ID,
		Username:  dto.Username,
		Email:     dto.Email,
		Password:  dto.Password,
		CreatedAt: dto.CreatedAt,
	}, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, dto := range r.db.Users {
		if dto.Username == username {
			return &domain.User{
				ID:        dto.ID,
				Username:  dto.Username,
				Email:     dto.Email,
				Password:  dto.Password,
				CreatedAt: dto.CreatedAt,
			}, nil
		}
	}

	return nil, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, dto := range r.db.Users {
		if dto.Email == email {
			return &domain.User{
				ID:        dto.ID,
				Username:  dto.Username,
				Email:     dto.Email,
				Password:  dto.Password,
				CreatedAt: dto.CreatedAt,
			}, nil
		}
	}

	return nil, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.db.Users[user.ID]
	if !exists {
		return nil, nil
	}

	dto := &database.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}

	r.db.Users[user.ID] = dto

	return user, nil
}

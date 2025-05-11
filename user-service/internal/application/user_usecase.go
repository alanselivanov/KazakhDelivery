package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"unicode/utf8"

	"user-service/internal/domain"
	"user-service/internal/infrastructure/database"
	"user-service/internal/infrastructure/persistence"

	"github.com/redis/go-redis/v9"
)

type UserUseCase struct {
	repo  persistence.UserRepository
	cache *database.RedisCache
}

func NewUserUseCase(repo persistence.UserRepository, cache *database.RedisCache) *UserUseCase {
	return &UserUseCase{
		repo:  repo,
		cache: cache,
	}
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, username, email, password string) (*domain.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("username, email and password are required")
	}

	if !isValidEmail(email) {
		return nil, errors.New("invalid email format")
	}

	if utf8.RuneCountInString(password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	existing, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	existing, err = uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	user, err := domain.NewUser(username, email, password)
	if err != nil {
		return nil, err
	}

	createdUser, err := uc.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (uc *UserUseCase) AuthenticateUser(ctx context.Context, username, password string) (*domain.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (uc *UserUseCase) GetUserProfile(ctx context.Context, userID string) (*domain.Profile, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	cacheKey := fmt.Sprintf("user_profile:%s", userID)
	var profile domain.Profile

	if uc.cache != nil {
		err := uc.cache.Get(ctx, cacheKey, &profile)
		if err == nil {
			log.Printf("Data retrieved from cache for user %s", userID)
			return &profile, nil
		} else if err != redis.Nil {
			log.Printf("Redis error: %v", err)
		}
	}

	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	log.Printf("Data retrieved from database for user %s", userID)
	userProfile := user.ToProfile()

	if uc.cache != nil {
		if err := uc.cache.Set(ctx, cacheKey, userProfile); err != nil {
			log.Printf("Failed to cache user profile: %v", err)
		}
	}

	return userProfile, nil
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, userID, username, email string) (*domain.Profile, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if username != "" && username != user.Username {
		existing, err := uc.repo.GetByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != userID {
			return nil, errors.New("username already exists")
		}
		user.Username = username
	}

	if email != "" && email != user.Email {
		if !isValidEmail(email) {
			return nil, errors.New("invalid email format")
		}
		existing, err := uc.repo.GetByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != userID {
			return nil, errors.New("email already exists")
		}
		user.Email = email
	}

	updatedUser, err := uc.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	if err := uc.InvalidateUserCache(ctx, userID); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	} else {
		log.Printf("Cache invalidated for user %s", userID)
	}

	return updatedUser.ToProfile(), nil
}

func (uc *UserUseCase) InvalidateUserCache(ctx context.Context, userID string) error {
	if uc.cache == nil {
		return nil
	}

	cacheKey := fmt.Sprintf("user_profile:%s", userID)

	return uc.cache.Delete(ctx, cacheKey)
}

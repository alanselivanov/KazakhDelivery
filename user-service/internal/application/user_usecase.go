package application

import (
	"context"
	"errors"
	"regexp"
	"unicode/utf8"

	"user-service/internal/domain"
	"user-service/internal/infrastructure/persistence"
)

type UserUseCase struct {
	repo persistence.UserRepository
}

func NewUserUseCase(repo persistence.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
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

	return uc.repo.Create(ctx, user)
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

	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user.ToProfile(), nil
}

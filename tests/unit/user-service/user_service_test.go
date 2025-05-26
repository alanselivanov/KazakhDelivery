package userservice_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id string) error
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) (*User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *User) (*User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type AuthService struct {
	userRepo UserRepository
}

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{
		userRepo: repo,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("authentication failed")
	}

	if user.Password != password {
		return nil, errors.New("authentication failed")
	}

	user.Password = ""
	return user, nil
}

func TestAuthService_Login_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	testUser := &User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByUsername", ctx, "testuser").Return(testUser, nil)

	resultUser, err := service.Login(ctx, "testuser", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.Equal(t, testUser.ID, resultUser.ID)
	assert.Equal(t, testUser.Username, resultUser.Username)
	assert.Equal(t, testUser.Email, resultUser.Email)
	assert.Empty(t, resultUser.Password)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	testUser := &User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByUsername", ctx, "testuser").Return(testUser, nil)

	resultUser, err := service.Login(ctx, "testuser", "wrong-password")

	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Equal(t, "authentication failed", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_EmptyCredentials(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	resultUser1, err1 := service.Login(ctx, "", "password123")
	assert.Error(t, err1)
	assert.Nil(t, resultUser1)
	assert.Equal(t, "username and password are required", err1.Error())

	resultUser2, err2 := service.Login(ctx, "testuser", "")
	assert.Error(t, err2)
	assert.Nil(t, resultUser2)
	assert.Equal(t, "username and password are required", err2.Error())

	mockRepo.AssertNotCalled(t, "FindByUsername")
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("FindByUsername", ctx, "nonexistent").Return(nil, errors.New("user not found"))

	resultUser, err := service.Login(ctx, "nonexistent", "anypassword")

	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Equal(t, "authentication failed", err.Error())
	mockRepo.AssertExpectations(t)
}

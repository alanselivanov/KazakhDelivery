package handlers

import (
	"context"
	"time"

	"proto/user"
	"user-service/internal/application"
)

type UserHandler struct {
	user.UnimplementedUserServiceServer
	userUseCase *application.UserUseCase
}

func NewUserHandler(userUseCase *application.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	domainUser, err := h.userUseCase.RegisterUser(ctx, req.User.Username, req.User.Email, req.User.Password)
	if err != nil {
		return nil, err
	}

	return &user.UserResponse{
		User: &user.User{
			Id:        domainUser.ID,
			Username:  domainUser.Username,
			Email:     domainUser.Email,
			CreatedAt: domainUser.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *UserHandler) AuthenticateUser(ctx context.Context, req *user.AuthRequest) (*user.AuthResponse, error) {
	domainUser, err := h.userUseCase.AuthenticateUser(ctx, req.Username, req.Password)
	if err != nil {
		return &user.AuthResponse{Success: false}, nil
	}

	token := "token_" + domainUser.ID

	return &user.AuthResponse{
		Token:   token,
		Success: true,
	}, nil
}

func (h *UserHandler) GetUserProfile(ctx context.Context, req *user.UserID) (*user.UserProfile, error) {
	profile, err := h.userUseCase.GetUserProfile(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &user.UserProfile{
		Id:       profile.ID,
		Username: profile.Username,
		Email:    profile.Email,
	}, nil
}

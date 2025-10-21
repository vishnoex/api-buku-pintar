package service

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"buku-pintar/pkg/oauth2"
	"context"
	"errors"
)

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) service.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) RegisterWithOAuth2(ctx context.Context, user *entity.User, provider oauth2.Provider) error {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}

	// Set default role and status if not provided
	if user.Role == "" {
		user.Role = entity.RoleReader
	}
	if user.Status == "" {
		user.Status = entity.StatusActive
	}

	// Create user in database
	return s.userRepo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *userService) UpdateUser(ctx context.Context, user *entity.User) error {
	// OAuth2 users don't have passwords, so we skip password handling
	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

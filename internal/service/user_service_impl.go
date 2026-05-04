package service

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"context"
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

func (s *userService) CreateUser(ctx context.Context, user *entity.User) error {
	if user.Role == "" {
		user.Role = entity.RoleReader
	}
	if user.Status == "" {
		user.Status = entity.StatusActive
	}
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

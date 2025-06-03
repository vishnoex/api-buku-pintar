package usecase

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"context"
	"errors"
)

type userUsecase struct {
	userRepo repository.UserRepository
	userSvc  service.UserService
}

// NewUserUsecase creates a new instance of UserUsecase
func NewUserUsecase(userRepo repository.UserRepository, userSvc service.UserService) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		userSvc:  userSvc,
	}
}

func (u *userUsecase) Register(ctx context.Context, user *entity.User) error {
	// Check if user already exists
	existingUser, err := u.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}

	// Use service for business logic
	return u.userSvc.Register(ctx, user)
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return u.userSvc.GetUserByID(ctx, id)
}

func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.userSvc.GetUserByEmail(ctx, email)
}

func (u *userUsecase) UpdateUser(ctx context.Context, user *entity.User) error {
	// Check if user exists
	existingUser, err := u.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	return u.userSvc.UpdateUser(ctx, user)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	existingUser, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	return u.userSvc.DeleteUser(ctx, id)
} 
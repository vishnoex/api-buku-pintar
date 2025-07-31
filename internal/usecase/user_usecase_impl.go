package usecase

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"context"
	"errors"
)

type userUsecase struct {
	userRepo     repository.UserRepository
	userService  service.UserService
}

// NewUserUsecase creates a new instance of UserUsecase
func NewUserUsecase(userRepo repository.UserRepository, userService service.UserService) UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		userService: userService,
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
	return u.userService.Register(ctx, user)
}

func (u *userUsecase) RegisterWithFirebase(ctx context.Context, user *entity.User, idToken string) error {
	return u.userService.RegisterWithFirebase(ctx, user, idToken)
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return u.userService.GetUserByID(ctx, id)
}

func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.userService.GetUserByEmail(ctx, email)
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

	return u.userService.UpdateUser(ctx, user)
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

	return u.userService.DeleteUser(ctx, id)
} 
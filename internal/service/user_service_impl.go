package service

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"buku-pintar/internal/domain/service"
	"buku-pintar/pkg/oauth2"
	"context"
	"errors"

	"firebase.google.com/go/v4/auth"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo     repository.UserRepository
	firebaseAuth *auth.Client
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository, firebaseAuth *auth.Client) service.UserService {
	return &userService{
		userRepo:     userRepo,
		firebaseAuth: firebaseAuth,
	}
}

func (s *userService) Register(ctx context.Context, user *entity.User) error {
	// Hash password
	if user.Password == nil {
		return errors.New("password is required")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedPasswordStr := string(hashedPassword)
	user.Password = &hashedPasswordStr

	return s.userRepo.Create(ctx, user)
}

func (s *userService) RegisterWithFirebase(ctx context.Context, user *entity.User, idToken string) error {
	// Verify the ID token
	token, err := s.firebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return errors.New("invalid ID token")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}

	// Set user ID from Firebase UID
	user.ID = token.UID

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
	// If password is being updated, hash it
	if user.Password != nil && *user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		hashedPasswordStr := string(hashedPassword)
		user.Password = &hashedPasswordStr
	}

	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

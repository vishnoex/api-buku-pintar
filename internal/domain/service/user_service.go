package service

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// UserService defines the interface for user business operations
type UserService interface {
	Register(ctx context.Context, user *entity.User) error
	RegisterWithFirebase(ctx context.Context, user *entity.User, idToken string) error
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id string) error
}

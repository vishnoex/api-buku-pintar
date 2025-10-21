package usecase

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/pkg/oauth2"
	"context"
)

type UserUsecase interface {
	RegisterWithOAuth2(ctx context.Context, user *entity.User, provider oauth2.Provider) error
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id string) error
} 
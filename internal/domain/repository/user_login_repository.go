package repository

import (
	"buku-pintar/internal/domain/entity"
	"context"
)

// UserLoginRepository defines the interface for user login operations
// Clean Architecture: Domain layer, no infrastructure dependencies
type UserLoginRepository interface {
	Create(ctx context.Context, userLogin *entity.UserLogin) error
	GetByID(ctx context.Context, id string) (*entity.UserLogin, error)
	GetByUserAndProvider(ctx context.Context, userID, loginProviderID string) (*entity.UserLogin, error)
	Update(ctx context.Context, userLogin *entity.UserLogin) error
	Delete(ctx context.Context, id string) error
	ListByUser(ctx context.Context, userID string) ([]*entity.UserLogin, error)
	ListActiveByUser(ctx context.Context, userID string) ([]*entity.UserLogin, error)
	CountByUser(ctx context.Context, userID string) (int64, error)
	CountActiveByUser(ctx context.Context, userID string) (int64, error)
}

package usecase

import (
	"buku-pintar/internal/delivery/http/response"
	"context"
)

type CategoryUsecase interface {
	ListCategory(ctx context.Context, limit, offset int) ([]*response.CategoryResponse, error)
	CountCategory(ctx context.Context) (int64, error)
}

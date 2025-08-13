package usecase

import (
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/repository"
	"context"
)

type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (uc *categoryUsecase) ListCategory(ctx context.Context, limit, offset int) ([]*response.CategoryResponse, error) {
	categories := []*response.CategoryResponse{}
	categoryEntities, err := uc.categoryRepo.ListActive(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, category := range categoryEntities {
		categories = append(categories, &response.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Slug:        category.Slug,
			Description: category.Description,
			Icon:        category.IconLink,
			OrderNumber: category.OrderNumber,
			IsActive:    category.IsActive,
		})
	}

	return categories, nil
}

// CountCategory implements CategoryUsecase.
func (uc *categoryUsecase) CountCategory(ctx context.Context) (int64, error) {
	return uc.categoryRepo.Count(ctx)
}

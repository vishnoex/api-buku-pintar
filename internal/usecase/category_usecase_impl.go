package usecase

import (
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"context"
)

type categoryUsecase struct {
	categoryService service.CategoryService
}

func NewCategoryUsecase(categoryService service.CategoryService) CategoryUsecase {
	return &categoryUsecase{
		categoryService: categoryService,
	}
}

func (uc *categoryUsecase) ListCategory(ctx context.Context, limit, offset int) ([]*response.CategoryResponse, error) {
	// Get categories from service layer (handles caching)
	categories, err := uc.categoryService.GetActiveCategoryList(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert entities to response DTOs
	var response []*response.CategoryResponse
	for _, category := range categories {
		categoryResponse := uc.convertEntityToResponse(category)
		response = append(response, &categoryResponse)
	}

	return response, nil
}

func (uc *categoryUsecase) GetCategoryByID(ctx context.Context, id string) (*response.CategoryResponse, error) {
	// Get category from service layer
	category, err := uc.categoryService.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, nil
	}

	// Convert entity to response DTO
	categoryResponse := uc.convertEntityToResponse(category)
	return &categoryResponse, nil
}

func (uc *categoryUsecase) CreateCategory(ctx context.Context, category *entity.Category) error {
	// Create category through service layer (handles cache invalidation)
	return uc.categoryService.CreateCategory(ctx, category)
}

func (uc *categoryUsecase) UpdateCategory(ctx context.Context, category *entity.Category) error {
	// Update category through service layer (handles cache invalidation)
	return uc.categoryService.UpdateCategory(ctx, category)
}

func (uc *categoryUsecase) DeleteCategory(ctx context.Context, id string) error {
	// Delete category through service layer (handles cache invalidation)
	return uc.categoryService.DeleteCategory(ctx, id)
}

func (uc *categoryUsecase) CountCategory(ctx context.Context) (int64, error) {
	// Get count from service layer (handles caching)
	return uc.categoryService.GetActiveCategoryCount(ctx)
}

func (uc *categoryUsecase) ListAllCategories(ctx context.Context, limit, offset int) ([]*response.CategoryResponse, error) {
	// Get all categories from service layer (handles caching)
	categories, err := uc.categoryService.GetCategoryList(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert entities to response DTOs
	var response []*response.CategoryResponse
	for _, category := range categories {
		categoryResponse := uc.convertEntityToResponse(category)
		response = append(response, &categoryResponse)
	}

	return response, nil
}

func (uc *categoryUsecase) CountAllCategories(ctx context.Context) (int64, error) {
	// Get total count from service layer (handles caching)
	return uc.categoryService.GetCategoryCount(ctx)
}

func (uc *categoryUsecase) ListCategoriesByParent(ctx context.Context, parentID string, limit, offset int) ([]*response.CategoryResponse, error) {
	// Get categories by parent from service layer
	categories, err := uc.categoryService.GetCategoriesByParent(ctx, parentID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert entities to response DTOs
	var response []*response.CategoryResponse
	for _, category := range categories {
		categoryResponse := uc.convertEntityToResponse(category)
		response = append(response, &categoryResponse)
	}

	return response, nil
}

func (uc *categoryUsecase) CountCategoriesByParent(ctx context.Context, parentID string) (int64, error) {
	// Get count by parent from service layer
	return uc.categoryService.GetCategoryCountByParent(ctx, parentID)
}

// convertEntityToResponse converts Category entity to CategoryResponse
func (uc *categoryUsecase) convertEntityToResponse(category *entity.Category) response.CategoryResponse {
	response := response.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Icon:        category.IconLink,
		OrderNumber: category.OrderNumber,
		IsActive:    category.IsActive,
	}
	
	if category.Description != nil {
		response.Description = category.Description
	}
	
	return response
}

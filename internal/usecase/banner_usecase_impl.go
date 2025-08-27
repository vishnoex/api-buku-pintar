package usecase

import (
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"context"
)

type bannerUsecase struct {
	bannerService service.BannerService
}

func NewBannerUsecase(bannerService service.BannerService) BannerUsecase {
	return &bannerUsecase{
		bannerService: bannerService,
	}
}

func (uc *bannerUsecase) ListBanner(ctx context.Context, limit, offset int) ([]*response.BannerResponse, error) {
	// Get banners from service layer (handles caching)
	banners, err := uc.bannerService.GetBannerList(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert entities to response DTOs
	var response []*response.BannerResponse
	for _, banner := range banners {
		bannerResponse := uc.convertEntityToResponse(banner)
		response = append(response, &bannerResponse)
	}

	return response, nil
}

func (uc *bannerUsecase) GetBannerByID(ctx context.Context, id string) (*response.BannerResponse, error) {
	// Get banner from service layer
	banner, err := uc.bannerService.GetBannerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if banner == nil {
		return nil, nil
	}

	// Convert entity to response DTO
	bannerResponse := uc.convertEntityToResponse(banner)
	return &bannerResponse, nil
}

func (uc *bannerUsecase) CreateBanner(ctx context.Context, banner *entity.Banner) error {
	// Create banner through service layer (handles cache invalidation)
	return uc.bannerService.CreateBanner(ctx, banner)
}

func (uc *bannerUsecase) UpdateBanner(ctx context.Context, banner *entity.Banner) error {
	// Update banner through service layer (handles cache invalidation)
	return uc.bannerService.UpdateBanner(ctx, banner)
}

func (uc *bannerUsecase) DeleteBanner(ctx context.Context, id string) error {
	// Delete banner through service layer (handles cache invalidation)
	return uc.bannerService.DeleteBanner(ctx, id)
}

func (uc *bannerUsecase) CountBanner(ctx context.Context) (int64, error) {
	// Get count from service layer (handles caching)
	return uc.bannerService.GetBannerCount(ctx)
}

func (uc *bannerUsecase) ListActiveBanner(ctx context.Context, limit, offset int) ([]*response.BannerResponse, error) {
	// Get active banners from service layer
	banners, err := uc.bannerService.GetActiveBannerList(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert entities to response DTOs
	var response []*response.BannerResponse
	for _, banner := range banners {
		bannerResponse := uc.convertEntityToResponse(banner)
		response = append(response, &bannerResponse)
	}

	return response, nil
}

func (uc *bannerUsecase) CountActiveBanner(ctx context.Context) (int64, error) {
	// Get active banner count from service layer
	return uc.bannerService.GetActiveBannerCount(ctx)
}

// convertEntityToResponse converts Banner entity to BannerResponse
func (uc *bannerUsecase) convertEntityToResponse(banner *entity.Banner) response.BannerResponse {
	response := response.BannerResponse{
		ID:    banner.ID,
		Title: banner.Title,
		Image: banner.ImageURL,
	}
	
	if banner.Link != nil {
		response.Link = *banner.Link
	}
	if banner.CTALabel != nil {
		response.CTALabel = *banner.CTALabel
	}
	if banner.BackgroundColor != nil {
		response.BackgroundColor = *banner.BackgroundColor
	}
	
	return response
}

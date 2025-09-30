package usecase

import (
	"buku-pintar/internal/constant"
	"buku-pintar/internal/delivery/http/response"
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/service"
	"context"
	"errors"
)

type ebookUsecase struct {
	ebookService service.EbookService
	ebookDiscountService service.EbookDiscountService
}

// NewEbookUsecase creates a new instance of EbookUsecase
func NewEbookUsecase(ebookService service.EbookService, ebookDiscountService service.EbookDiscountService) EbookUsecase {
	return &ebookUsecase{
		ebookService: ebookService,
		ebookDiscountService: ebookDiscountService,
	}
}

func (u *ebookUsecase) CreateEbook(ctx context.Context, ebook *entity.Ebook) error {
	// Validate required fields
	if ebook.Title == "" {
		return errors.New("title is required")
	}
	if ebook.AuthorID == "" {
		return errors.New(constant.ERR_AUTHOR_ID_REQUIRED)
	}
	if ebook.CategoryID == "" {
		return errors.New(constant.ERR_CATEGORY_ID_REQUIRED)
	}
	if ebook.Slug == "" {
		return errors.New("slug is required")
	}

	// Check if ebook with same slug already exists
	existingEbook, err := u.ebookService.GetEbookBySlug(ctx, ebook.Slug)
	if err != nil {
		return err
	}
	if existingEbook != nil {
		return errors.New("ebook with this slug already exists")
	}

	return u.ebookService.CreateEbook(ctx, ebook)
}

func (u *ebookUsecase) GetEbookByID(ctx context.Context, id string) (*entity.Ebook, error) {
	if id == "" {
		return nil, errors.New(constant.ERR_ID_REQUIRED)
	}

	return u.ebookService.GetEbookByID(ctx, id)
}

func (u *ebookUsecase) GetEbookBySlug(ctx context.Context, slug string) (*response.EbookResponse, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}

	ebook, err := u.ebookService.GetEbookBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if ebook == nil {
		return nil, nil
	}

	res := response.ParseEbookResponse(ebook)
	discount, _ := u.ebookDiscountService.GetActiveDiscountByEbookID(ctx, ebook.ID)
	if discount != nil {
		res.Discount = response.ParseDiscountResponse(discount)
	}

	return res, nil
}

func (u *ebookUsecase) UpdateEbook(ctx context.Context, ebook *entity.Ebook) error {
	if ebook.ID == "" {
		return errors.New(constant.ERR_ID_REQUIRED)
	}

	// Check if ebook exists
	existingEbook, err := u.ebookService.GetEbookByID(ctx, ebook.ID)
	if err != nil {
		return err
	}
	if existingEbook == nil {
		return errors.New("ebook not found")
	}

	// If slug is being updated, check for uniqueness
	if ebook.Slug != "" && ebook.Slug != existingEbook.Slug {
		existingBySlug, err := u.ebookService.GetEbookBySlug(ctx, ebook.Slug)
		if err != nil {
			return err
		}
		if existingBySlug != nil {
			return errors.New("ebook with this slug already exists")
		}
	}

	return u.ebookService.UpdateEbook(ctx, ebook)
}

func (u *ebookUsecase) DeleteEbook(ctx context.Context, id string) error {
	if id == "" {
		return errors.New(constant.ERR_ID_REQUIRED)
	}

	// Check if ebook exists
	existingEbook, err := u.ebookService.GetEbookByID(ctx, id)
	if err != nil {
		return err
	}
	if existingEbook == nil {
		return errors.New("ebook not found")
	}

	return u.ebookService.DeleteEbook(ctx, id)
}

func (u *ebookUsecase) ListEbooks(ctx context.Context, limit, offset int) ([]*response.EbookListResponse, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0
	}

	ebooks := []*response.EbookListResponse{}
	ebookList, err := u.ebookService.GetEbookList(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, ebook := range ebookList {
		ebooks = append(ebooks, response.ParseEbookListResponse(ebook))
	}

	return ebooks, nil
}

func (u *ebookUsecase) ListEbooksByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error) {
	if categoryID == "" {
		return nil, errors.New(constant.ERR_CATEGORY_ID_REQUIRED)
	}

	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0
	}

	return u.ebookService.GetEbookListByCategory(ctx, categoryID, limit, offset)
}

func (u *ebookUsecase) ListEbooksByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error) {
	if authorID == "" {
		return nil, errors.New(constant.ERR_AUTHOR_ID_REQUIRED)
	}

	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0
	}

	return u.ebookService.GetEbookListByAuthor(ctx, authorID, limit, offset)
}

func (u *ebookUsecase) CountEbooks(ctx context.Context) (int64, error) {
	return u.ebookService.GetEbookCount(ctx)
}

func (u *ebookUsecase) CountEbooksByCategory(ctx context.Context, categoryID string) (int64, error) {
	if categoryID == "" {
		return 0, errors.New(constant.ERR_CATEGORY_ID_REQUIRED)
	}

	return u.ebookService.GetEbookCountByCategory(ctx, categoryID)
}

func (u *ebookUsecase) CountEbooksByAuthor(ctx context.Context, authorID string) (int64, error) {
	if authorID == "" {
		return 0, errors.New(constant.ERR_AUTHOR_ID_REQUIRED)
	}

	return u.ebookService.GetEbookCountByAuthor(ctx, authorID)
}

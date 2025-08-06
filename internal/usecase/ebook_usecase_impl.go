package usecase

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"errors"
)

type ebookUsecase struct {
	ebookRepo repository.EbookRepository
}

// NewEbookUsecase creates a new instance of EbookUsecase
func NewEbookUsecase(ebookRepo repository.EbookRepository) EbookUsecase {
	return &ebookUsecase{
		ebookRepo: ebookRepo,
	}
}

func (u *ebookUsecase) CreateEbook(ctx context.Context, ebook *entity.Ebook) error {
	// Validate required fields
	if ebook.Title == "" {
		return errors.New("title is required")
	}
	if ebook.AuthorID == "" {
		return errors.New("author_id is required")
	}
	if ebook.CategoryID == "" {
		return errors.New("category_id is required")
	}
	if ebook.Slug == "" {
		return errors.New("slug is required")
	}

	// Check if ebook with same slug already exists
	existingEbook, err := u.ebookRepo.GetBySlug(ctx, ebook.Slug)
	if err != nil {
		return err
	}
	if existingEbook != nil {
		return errors.New("ebook with this slug already exists")
	}

	return u.ebookRepo.Create(ctx, ebook)
}

func (u *ebookUsecase) GetEbookByID(ctx context.Context, id string) (*entity.Ebook, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return u.ebookRepo.GetByID(ctx, id)
}

func (u *ebookUsecase) GetEbookBySlug(ctx context.Context, slug string) (*entity.Ebook, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}

	return u.ebookRepo.GetBySlug(ctx, slug)
}

func (u *ebookUsecase) UpdateEbook(ctx context.Context, ebook *entity.Ebook) error {
	if ebook.ID == "" {
		return errors.New("id is required")
	}

	// Check if ebook exists
	existingEbook, err := u.ebookRepo.GetByID(ctx, ebook.ID)
	if err != nil {
		return err
	}
	if existingEbook == nil {
		return errors.New("ebook not found")
	}

	// If slug is being updated, check for uniqueness
	if ebook.Slug != "" && ebook.Slug != existingEbook.Slug {
		existingBySlug, err := u.ebookRepo.GetBySlug(ctx, ebook.Slug)
		if err != nil {
			return err
		}
		if existingBySlug != nil {
			return errors.New("ebook with this slug already exists")
		}
	}

	return u.ebookRepo.Update(ctx, ebook)
}

func (u *ebookUsecase) DeleteEbook(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	// Check if ebook exists
	existingEbook, err := u.ebookRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingEbook == nil {
		return errors.New("ebook not found")
	}

	return u.ebookRepo.Delete(ctx, id)
}

func (u *ebookUsecase) ListEbooks(ctx context.Context, limit, offset int) ([]*entity.Ebook, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0
	}

	return u.ebookRepo.List(ctx, limit, offset)
}

func (u *ebookUsecase) ListEbooksByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Ebook, error) {
	if categoryID == "" {
		return nil, errors.New("category_id is required")
	}

	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0
	}

	return u.ebookRepo.ListByCategory(ctx, categoryID, limit, offset)
}

func (u *ebookUsecase) ListEbooksByAuthor(ctx context.Context, authorID string, limit, offset int) ([]*entity.Ebook, error) {
	if authorID == "" {
		return nil, errors.New("author_id is required")
	}

	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // default limit
	}
	if offset < 0 {
		offset = 0
	}

	return u.ebookRepo.ListByAuthor(ctx, authorID, limit, offset)
}

func (u *ebookUsecase) CountEbooks(ctx context.Context) (int64, error) {
	return u.ebookRepo.Count(ctx)
}

func (u *ebookUsecase) CountEbooksByCategory(ctx context.Context, categoryID string) (int64, error) {
	if categoryID == "" {
		return 0, errors.New("category_id is required")
	}

	return u.ebookRepo.CountByCategory(ctx, categoryID)
}

func (u *ebookUsecase) CountEbooksByAuthor(ctx context.Context, authorID string) (int64, error) {
	if authorID == "" {
		return 0, errors.New("author_id is required")
	}

	return u.ebookRepo.CountByAuthor(ctx, authorID)
}

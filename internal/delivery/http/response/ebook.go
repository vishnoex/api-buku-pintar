package response

import (
	"buku-pintar/internal/domain/entity"
	"strconv"
)

type EbookResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Synopsis    string `json:"synopsis"`
	Slug        string `json:"slug"`
	CoverImage  string `json:"cover_image"`
	Status      string `json:"status"`
	Price       int    `json:"price"`
	Language    string `json:"language"`
	Duration    int    `json:"duration"`
	Filesize    int64  `json:"filesize"`
	Format      string `json:"format"`
	PageCount   int16  `json:"page_count"`
	PreviewPage int16  `json:"preview_page"`
	URL         string `json:"url"`
	PublishedAt string `json:"published_at"`

	Author          AuthorResponse               `json:"author"`
	Category        *CategoryResponse            `json:"category"`
	Discount        *EbookDiscountResponse       `json:"discount"`
	TableOfContents []*TableOfContentResponse    `json:"table_of_contents"`
	Summary         *EbookSummaryResponse        `json:"summary"`
	PremiumSummary  *EbookPremiumSummaryResponse `json:"premium_summary"`
}

type EbookListResponse struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	CoverImage string `json:"cover_image"`
	Status     string `json:"status"`
	Price      int    `json:"price"`
	Discount   int    `json:"discount"`
}

func ParseEbookListResponse(ebook *entity.EbookList) *EbookListResponse {
	return &EbookListResponse{
		ID:         ebook.ID,
		Title:      ebook.Title,
		Slug:       ebook.Slug,
		CoverImage: ebook.CoverImage,
		Price:      ebook.Price,
		Discount:   0,
	}
}

func ParseEbookResponse(ebook *entity.EbookDetail) *EbookResponse {
	// Format the published date as string
	var publishedAtStr string
	if ebook.PublishedAt != nil {
		publishedAtStr = ebook.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return &EbookResponse{
		ID:          ebook.ID,
		Title:       ebook.Title,
		Synopsis:    ebook.Synopsis,
		Slug:        ebook.Slug,
		CoverImage:  ebook.CoverImage,
		Status:      "", // This field doesn't exist in the entity, will be empty
		Price:       ebook.Price,
		Language:    ebook.Language,
		Duration:    ebook.Duration,
		Filesize:    ebook.Filesize,
		Format:      string(ebook.Format), // Convert EbookFormat to string
		PageCount:   ebook.PageCount,
		PreviewPage: ebook.PreviewPage,
		URL:         ebook.URL,
		PublishedAt: publishedAtStr,
		// Related entities will be empty for now since they're not included in the basic entity
		Author: AuthorResponse{
			ID:     ebook.AuthorID,
			Name:   ebook.AuthorName,
			Avatar: ebook.AuthorAvatar,
		},
		Category:        nil,
		Discount:        nil,
		TableOfContents: []*TableOfContentResponse{},
		Summary: &EbookSummaryResponse{
			ID:          *ebook.SummaryID,
			EbookID:     ebook.ID,
			EbookTitle:  ebook.Title,
			Slug:        ebook.Slug,
			Description: *ebook.SummaryContent,
			URL:         ebook.URL,
			AudioURL:    *ebook.SummaryAudioURL,
			Duration:    strconv.Itoa(*ebook.SummaryDuration),
		},
		PremiumSummary: nil,
	}
}

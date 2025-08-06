package response

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
	Category        CategoryResponse             `json:"category"`
	Discount        *EbookDiscountResponse       `json:"discount"`
	TableOfContents []*TableOfContentResponse    `json:"table_of_contents"`
	Summary         *EbookSummaryResponse        `json:"summary"`
	PremiumSummary  *EbookPremiumSummaryResponse `json:"premium_summary"`
}

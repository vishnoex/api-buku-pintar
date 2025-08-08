package entity

import "time"

// EbookFormat represents the format of an ebook
type EbookFormat string

const (
	FormatPDF  EbookFormat = "pdf"
	FormatEPUB EbookFormat = "epub"
	FormatMOBI EbookFormat = "mobi"
)

// Ebook represents an ebook in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
type Ebook struct {
	ID              string      `db:"id" json:"id"`
	AuthorID        string      `db:"author_id" json:"author_id"`
	Title           string      `db:"title" json:"title"`
	Synopsis        string      `db:"synopsis" json:"synopsis"`
	Slug            string      `db:"slug" json:"slug"`
	CoverImage      string      `db:"cover_image" json:"cover_image"`
	CategoryID      string      `db:"category_id" json:"category_id"`
	ContentStatusID string      `db:"content_status_id" json:"content_status_id"`
	Price           int         `db:"price" json:"price"`
	Language        string      `db:"language" json:"language"`
	Duration        int         `db:"duration" json:"duration"`
	Filesize        int64       `db:"filesize" json:"filesize"`
	Format          EbookFormat `db:"format" json:"format"`
	PageCount       int16       `db:"page_count" json:"page_count"`
	PreviewPage     int16       `db:"preview_page" json:"preview_page"`
	URL             string      `db:"url" json:"url"`
	PublishedAt     *time.Time  `db:"published_at" json:"published_at"`
	CreatedAt       time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at" json:"updated_at"`
}

type EbookList struct {
	ID         string `db:"id"`
	Title      string `db:"title"`
	Slug       string `db:"slug"`
	CoverImage string `db:"cover_image"`
	Price      int    `db:"price"`
	Discount   *int   `db:"discount"`
}

type EbookDetail struct {
	ID          string      `db:"id"`
	Title       string      `db:"title"`
	Synopsis    string      `db:"synopsis"`
	Slug        string      `db:"slug"`
	CoverImage  string      `db:"cover_image"`
	Price       int         `db:"price"`
	Language    string      `db:"language"`
	Duration    int         `db:"duration"`
	Filesize    int64       `db:"filesize"`
	Format      EbookFormat `db:"format"`
	PageCount   int16       `db:"page_count"`
	PreviewPage int16       `db:"preview_page"`
	URL         string      `db:"url"`
	PublishedAt *time.Time  `db:"published_at"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at"`

	AuthorID     string  `db:"author_id"`
	AuthorName   string  `db:"author_name"`
	AuthorAvatar *string `db:"author_avatar"`

	ContentStatus *string `db:"content_status"`
}

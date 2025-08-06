package entity

import "time"

// TableOfContent represents a table of content entry in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
type TableOfContent struct {
	ID          string    `db:"id" json:"id"`
	EbookID     string    `db:"ebook_id" json:"ebook_id"`
	Title       string    `db:"title" json:"title"`
	PageNumber  int16     `db:"page_number" json:"page_number"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

package entity

import "time"

// EbookSummary represents a summary for an ebook
type EbookSummary struct {
	ID          string    `db:"id" json:"id"`
	EbookID     string    `db:"ebook_id" json:"ebook_id"`
	Description string    `db:"description" json:"description"`
	URL         string    `db:"url" json:"url"`
	AudioURL    string    `db:"audio_url" json:"audio_url"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type EbookSummaryList struct {
	ID          string    `db:"id" json:"id"`
	EbookID     string    `db:"ebook_id" json:"ebook_id"`
	EbookTitle  string    `db:"ebook_title" json:"ebook_title"`
	Slug        string    `db:"slug" json:"slug"`
	Description string    `db:"description" json:"description"`
	URL         string    `db:"url" json:"url"`
	AudioURL    string    `db:"audio_url" json:"audio_url"`
	Duration    int	      `db:"duration" json:"duration"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

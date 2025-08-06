package entity

import "time"

// EbookPremiumSummary represents a premium summary for an ebook
type EbookPremiumSummary struct {
	ID          string    `db:"id" json:"id"`
	EbookID     string    `db:"ebook_id" json:"ebook_id"`
	Description string    `db:"description" json:"description"`
	URL         string    `db:"url" json:"url"`
	AudioURL    string    `db:"audio_url" json:"audio_url"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

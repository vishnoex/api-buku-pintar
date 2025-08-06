package entity

import "time"

// EbookDiscount represents a discount for an ebook
type EbookDiscount struct {
	ID            string    `db:"id" json:"id"`
	EbookID       string    `db:"ebook_id" json:"ebook_id"`
	DiscountPrice int       `db:"discount_price" json:"discount_price"`
	StartedAt     time.Time `db:"started_at" json:"started_at"`
	EndedAt       time.Time `db:"ended_at" json:"ended_at"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

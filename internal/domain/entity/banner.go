package entity

import "time"

// Banner represents a banner in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
type Banner struct {
	ID              string     `db:"id" json:"id"`
	Title           string     `db:"title" json:"title"`
	ImageURL        string     `db:"image_url" json:"image_url"`
	Link            *string    `db:"link" json:"link"`
	CTALabel        *string    `db:"cta_label" json:"cta_label"`
	BackgroundColor *string    `db:"background_color" json:"background_color"`
	IsActive        bool       `db:"is_active" json:"is_active"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

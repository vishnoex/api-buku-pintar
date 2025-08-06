package entity

import "time"

// Article represents an article in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
type Article struct {
	ID              string     `db:"id" json:"id"`
	AuthorID        string     `db:"author_id" json:"author_id"`
	Title           string     `db:"title" json:"title"`
	Content         string     `db:"content" json:"content"`
	Slug            string     `db:"slug" json:"slug"`
	Excerpt         string     `db:"excerpt" json:"excerpt"`
	CoverImage      string     `db:"cover_image" json:"cover_image"`
	CategoryID      string     `db:"category_id" json:"category_id"`
	ContentStatusID string     `db:"content_status_id" json:"content_status_id"`
	ReadingTime     int16      `db:"reading_time" json:"reading_time"`
	PublishedAt     *time.Time `db:"published_at" json:"published_at"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

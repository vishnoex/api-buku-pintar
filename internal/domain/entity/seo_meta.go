package entity

import "time"

// SeoEntityType represents the type of entity for SEO metadata
type SeoEntityType string

const (
	SeoEntityBanner      SeoEntityType = "banner"
	SeoEntityCategory    SeoEntityType = "category"
	SeoEntityEbook       SeoEntityType = "ebook"
	SeoEntityArticle     SeoEntityType = "article"
	SeoEntityInspiration SeoEntityType = "inspiration"
)

// SeoMeta represents SEO metadata in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
type SeoMeta struct {
	ID          string        `db:"id" json:"id"`
	Title       string        `db:"title" json:"title"`
	Description *string       `db:"description" json:"description"`
	Keywords    *string       `db:"keywords" json:"keywords"`
	Entity      SeoEntityType `db:"entity" json:"entity"`
	EntityID    string        `db:"entity_id" json:"entity_id"`
	CreatedAt   time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at" json:"updated_at"`
}

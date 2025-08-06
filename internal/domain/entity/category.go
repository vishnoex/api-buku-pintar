package entity

import "time"

// Category represents a category in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
type Category struct {
	ID          string     `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description"`
	IconLink    string     `db:"icon_link" json:"icon_link"`
	ParentID    *string    `db:"parent_id" json:"parent_id"`
	OrderNumber int8       `db:"order_number" json:"order_number"`
	IsActive    bool       `db:"is_active" json:"is_active"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

package entity

import "time"

// ContentStatus represents a content status in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
type ContentStatus struct {
	ID        string     `db:"id" json:"id"`
	Name      *string    `db:"name" json:"name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

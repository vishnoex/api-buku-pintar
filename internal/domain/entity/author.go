package entity

import "time"

// Author represents an author in the system
// Clean Architecture: Entity layer, no dependencies on infrastructure
type Author struct {
	ID        string     `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Avatar    *string    `db:"avatar" json:"avatar"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

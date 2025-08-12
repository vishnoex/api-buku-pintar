package response

type CategoryResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
	Icon        string  `json:"icon"`
	ParentID    *string `json:"parent_id,omitempty"`
	OrderNumber int8    `json:"order_number"`
	IsActive    bool    `json:"is_active"`
}

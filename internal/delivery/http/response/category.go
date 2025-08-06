package response

type CategoryResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IconLink    string  `json:"icon_link"`
	ParentID    *string `json:"parent_id"`
	OrderNumber int8    `json:"order_number"`
	IsActive    bool    `json:"is_active"`
}

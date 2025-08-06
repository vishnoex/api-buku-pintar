package response

type EbookDiscountResponse struct {
	ID            string `json:"id"`
	EbookID       string `json:"ebook_id"`
	DiscountPrice int    `json:"discount_price"`
	StartedAt     string `json:"started_at"`
	EndedAt       string `json:"ended_at"`
}

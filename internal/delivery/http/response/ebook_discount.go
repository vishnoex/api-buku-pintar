package response

import "buku-pintar/internal/domain/entity"

type EbookDiscountResponse struct {
	ID            string `json:"id"`
	EbookID       string `json:"ebook_id"`
	DiscountPrice int    `json:"discount_price"`
	StartedAt     string `json:"started_at"`
	EndedAt       string `json:"ended_at"`
}


func ParseDiscountResponse(discount *entity.EbookDiscount) *EbookDiscountResponse {
	return &EbookDiscountResponse{
		ID:            discount.ID,
		EbookID:       discount.EbookID,
		DiscountPrice: discount.DiscountPrice,
		StartedAt:     discount.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
		EndedAt:       discount.EndedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

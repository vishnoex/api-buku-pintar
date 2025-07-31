package xendit

import "time"

// Promotion contains data from Xendit's API response of promotion-related request.
// For more details see https://xendit.github.io/apireference/?bash#create-promotion.
type Promotion struct {
	ID                string     `json:"id"`
	BusinessID        string     `json:"business_id"`
	Status            string     `json:"status"`
	ReferenceID       string     `json:"reference_id"`
	Description       string     `json:"description"`
	PromoCode         string     `json:"promo_code"`
	BinList           []string   `json:"bin_list"`
	ChannelCode       string     `json:"channel_code"`
	DiscountPercent   float64    `json:"discount_percent"`
	DiscountAmount    float64    `json:"discount_amount"`
	Currency          string     `json:"currency"`
	StartTime         *time.Time `json:"start_time"`
	EndTime           *time.Time `json:"end_time"`
	MinOriginalAmount float64    `json:"min_original_amount"`
	MaxDiscountAmount float64    `json:"max_discount_amount"`
}

// PromotionDeletion contains data from Xendit's API response of delete promotion request.
// For more details see https://xendit.github.io/apireference/?bash#create-promotion.
type PromotionDeletion struct {
	ID        string `json:"id"`
	IsDeleted bool   `json:"is_deleted"`
}

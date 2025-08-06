package response

type EbookPremiumSummaryResponse struct {
	ID          string `json:"id"`
	EbookID     string `json:"ebook_id"`
	Description string `json:"description"`
	URL         string `json:"url"`
	AudioURL    string `json:"audio_url"`
}

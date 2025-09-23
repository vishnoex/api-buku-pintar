package response

type EbookSummaryResponse struct {
	ID          string `json:"id"`
	EbookID     string `json:"ebook_id"`
	EbookTitle  string `json:"ebook_title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	URL         string `json:"url"`
	AudioURL    string `json:"audio_url"`
	Duration    string `json:"duration"`
}

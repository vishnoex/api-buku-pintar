package response

type TableOfContentResponse struct {
	ID         string `json:"id"`
	EbookID    string `json:"ebook_id"`
	Title      string `json:"title"`
	PageNumber int16  `json:"page_number"`
}

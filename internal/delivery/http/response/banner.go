package response

type BannerResponse struct {
    ID string `json:"id"`
    Title string `json:"title"`
    Description string `json:"description"`
    Image string `json:"image"`
    Link string `json:"link"`
    CTALabel string `json:"cta_label"`
    BackgroundColor string `json:"background_color"`
}

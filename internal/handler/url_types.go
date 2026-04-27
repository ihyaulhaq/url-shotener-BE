package handler

type CreateURLRequest struct {
	OriginalUrl string `json:"original_url" validate:"required,url"`
}

type CreateURLResponse struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

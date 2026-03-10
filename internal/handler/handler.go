package handler

import (
	"net/http"

	"github.com/ihyaulhaq/url-shotener-BE/internal/service"
	"github.com/ihyaulhaq/url-shotener-BE/pkg/utils"
)

type Handler struct {
	urlService *service.UrlService
	baseURL    string
}

func New(urlService *service.UrlService, baseURL string) *Handler {
	return &Handler{
		urlService: urlService,
		baseURL:    baseURL,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.heatlh)

	mux.HandleFunc("POST /api/urls/shorten", h.handleShorteningUrl)
	mux.HandleFunc("GEt /api/urls/{shortUrl}", h.handleRedirectUrl)

	return mux
}

func (h *Handler) heatlh(w http.ResponseWriter, r *http.Request) {
	utils.ResponseWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})

}

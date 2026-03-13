package handler

import (
	"net/http"

	"github.com/ihyaulhaq/url-shotener-BE/internal/service"
	"github.com/ihyaulhaq/url-shotener-BE/pkg/utils"
)

type Handler struct {
	userService *service.UserService
	urlService  *service.UrlService
	baseURL     string
}

func New(
	urlService *service.UrlService,
	userService *service.UserService,
	baseURL string,
) *Handler {
	return &Handler{
		urlService:  urlService,
		userService: userService,
		baseURL:     baseURL,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.heatlh)

	mux.HandleFunc("POST /api/urls/shorten", h.handleShorteningUrl)
	mux.HandleFunc("GET /{shortUrl}", h.handleRedirectUrl)

	mux.HandleFunc("POST /login", h.hanldleUserLogin)

	return mux
}

func (h *Handler) heatlh(w http.ResponseWriter, r *http.Request) {
	utils.ResponseWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})

}

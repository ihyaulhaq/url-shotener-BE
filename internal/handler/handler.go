package handler

import (
	"net/http"

	"github.com/ihyaulhaq/url-shotener-BE/pkg/utils"
)

type Handler struct {
}

func New() *Handler {

	return &Handler{}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.heatlh)

	mux.HandleFunc("POST /api/urls/shorten", h.handleShorteningUrl)
	mux.HandleFunc("GEt /api/urls/{shortUrl}", h.handleRedirectUrl)

	return mux
}

func (h *Handler) heatlh(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, map[string]string{"status": "ok"})

}

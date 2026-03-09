package handler

import "net/http"

type Handler struct {
}

func New() *Handler {

	return &Handler{}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	return mux
}

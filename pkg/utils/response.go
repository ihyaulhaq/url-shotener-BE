package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Envelope struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
	Meta  *Meta     `json:"meta,omitempty"`
}

type APIError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Meta struct {
	TotalItems  int `json:"total_items,omitempty"`
	TotalPages  int `json:"total_pages,omitempty"`
	CurrentPage int `json:"current_page,omitempty"`
}

func ResponseWithJSON(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, Envelope{
		Data: data,
	})
}

func ResponseWithPagination(w http.ResponseWriter, status int, data any, meta Meta) {
	writeJSON(w, status, Envelope{
		Data: data,
		Meta: &meta,
	})
}
func writeJSON(w http.ResponseWriter, status int, envelope Envelope) {
	// Marshal first
	b, err := json.Marshal(envelope)
	if err != nil {
		slog.Error("failed to marshal JSON response", "error", err)
		http.Error(w,
			`{"error":{"code":500,"message":"internal server error"}}`,
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		slog.Error("failed to write response body", "error", err)
	}
}

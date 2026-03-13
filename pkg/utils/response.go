package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Envelope struct {
	Data      any       `json:"data,omitempty"`
	Error     *APIError `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ResponseWithJSON(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, Envelope{
		Data:      data,
		Timestamp: time.Now().UTC(),
	})
}

func ResponseWithError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, Envelope{
		Error: &APIError{
			Code:    status,
			Message: msg,
		},
		Timestamp: time.Now().UTC(),
	})
}

func writeJSON(w http.ResponseWriter, status int, envelope Envelope) {
	// Marshal first
	b, err := json.Marshal(envelope)
	if err != nil {
		slog.Error("failed to marshal JSON response", "error", err)
		http.Error(w, `{"error":{"code":500,"message":"internal server error"}}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(b); err != nil {
		slog.Error("failed to write response body", "error", err)
	}
}

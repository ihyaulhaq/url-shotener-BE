package utils

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Data  any    `json:"data"`
	Error string `json:"error,omitempty"`
}

func ResponseWithJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(
		Envelope{Data: data},
	)
}

func ResponseWithError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(
		Envelope{Error: msg},
	)
}

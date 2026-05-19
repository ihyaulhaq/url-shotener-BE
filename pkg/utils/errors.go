package utils

import (
	"log/slog"
	"net/http"
)

const (
	ErrTypeInvalidRequest = "invalid_request"
	ErrTypeValidation     = "validation_error"
	ErrTypeUnauthorized   = "unauthorized"
	ErrTypeForbidden      = "forbidden"
	ErrTypeNotFound       = "not_found"
	ErrTypeConflict       = "conflict"
	ErrTypeInternal       = "internal_error"
)

type ResponseError struct {
	Status  int
	Type    string
	Message string
	Cause   error
	Details any
}

func (e ResponseError) Write(w http.ResponseWriter) {
	if e.Cause != nil {
		slog.Error("request failed", "type", e.Type, "status", e.Status, "error", e.Cause)
	}
	writeJSON(w, e.Status, Envelope{
		Error: &APIError{Type: e.Type, Message: e.Message, Details: e.Details},
	})
}

// ErrInvalidRequest — malformed / unparseable JSON body.
func ErrInvalidRequest(cause error) ResponseError {
	return ResponseError{
		Status:  http.StatusBadRequest,
		Type:    ErrTypeInvalidRequest,
		Message: "request body is not valid JSON",
		Cause:   cause,
	}
}

// ErrBadRequest — valid JSON but failed field-level validation.
func ErrBadRequest(msg string) ResponseError {
	return ResponseError{
		Status:  http.StatusBadRequest,
		Type:    ErrTypeValidation,
		Message: msg,
	}
}

// ErrUnauthorized — missing or invalid credentials.
func ErrUnauthorized(msg string) ResponseError {
	return ResponseError{
		Status:  http.StatusUnauthorized,
		Type:    ErrTypeUnauthorized,
		Message: msg,
	}
}

// ErrForbidden — authenticated but not allowed.
func ErrForbidden(msg string) ResponseError {
	return ResponseError{
		Status:  http.StatusForbidden,
		Type:    ErrTypeForbidden,
		Message: msg,
	}
}

// ErrNotFound — resource doesn't exist.
func ErrNotFound(resource string) ResponseError {
	return ResponseError{
		Status:  http.StatusNotFound,
		Type:    ErrTypeNotFound,
		Message: resource + " not found",
	}
}

// ErrConflict — e.g. duplicate email.
func ErrConflict(msg string) ResponseError {
	return ResponseError{
		Status:  http.StatusConflict,
		Type:    ErrTypeConflict,
		Message: msg,
	}
}

// ErrInternal — catch-all server fault, always logs cause.
func ErrInternal(cause error) ResponseError {
	return ResponseError{
		Status:  http.StatusInternalServerError,
		Type:    ErrTypeInternal,
		Message: "something went wrong",
		Cause:   cause,
	}
}

package transport

import (
	"context"
	"encoding/json"
	"net/http"
)

// HTTPHandler defines a standard HTTP handler interface.
type HTTPHandler interface {
	RegisterRoutes(r Router)
}

// Router abstracts HTTP routing (works with Gin, Fiber, stdlib).
type Router interface {
	GET(path string, handler http.HandlerFunc)
	POST(path string, handler http.HandlerFunc)
	PUT(path string, handler http.HandlerFunc)
	PATCH(path string, handler http.HandlerFunc)
	DELETE(path string, handler http.HandlerFunc)
	Group(prefix string) Router
	Use(middleware ...func(http.Handler) http.Handler)
}

// Response is a standard API response envelope.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Error is a standard API error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta contains pagination or request metadata.
type Meta struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	Total    int `json:"total,omitempty"`
}

// WriteJSON writes a JSON response.
func WriteJSON(w http.ResponseWriter, status int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonEncoder := func(v interface{}) error {
		enc := json.NewEncoder(w)
		return enc.Encode(v)
	}
	jsonEncoder(resp)
}

// BadRequest returns a 400 response.
func BadRequest(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusBadRequest, Response{
		Success: false,
		Error: &Error{
			Code:    "BAD_REQUEST",
			Message: message,
		},
	})
}

// NotFound returns a 404 response.
func NotFound(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusNotFound, Response{
		Success: false,
		Error: &Error{
			Code:    "NOT_FOUND",
			Message: message,
		},
	})
}

// InternalError returns a 500 response.
func InternalError(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusInternalServerError, Response{
		Success: false,
		Error: &Error{
			Code:    "INTERNAL_ERROR",
			Message: message,
		},
	})
}

// OK returns a 200 response.
func OK(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created returns a 201 response.
func Created(w http.ResponseWriter, data interface{}) {
	WriteJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

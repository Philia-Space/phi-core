package transport

import (
	"errors"
	"net/http"

	phierr "github.com/philiaspace/phi-utils/errors"
)

// errorCodeToStatus maps domain error codes to HTTP status codes.
//
// This is the single source of truth for HTTP status mapping across all
// Philia Space services. Any service that wants standardized error responses
// should route errors through FromError below.
var errorCodeToStatus = map[phierr.ErrorCode]int{
	phierr.ErrNotFound:       http.StatusNotFound,
	phierr.ErrAlreadyExists:  http.StatusConflict,
	phierr.ErrInvalidInput:   http.StatusBadRequest,
	phierr.ErrUnauthorized:   http.StatusUnauthorized,
	phierr.ErrForbidden:      http.StatusForbidden,
	phierr.ErrConflict:       http.StatusConflict,
	phierr.ErrInternal:       http.StatusInternalServerError,
	phierr.ErrNotImplemented: http.StatusNotImplemented,
}

// FromError writes an HTTP response derived from a domain error.
//
// If the error is a *phierr.DomainError, the response uses the mapped HTTP
// status and exposes the error's code+message+details. If the error is not a
// DomainError (e.g. an unexpected runtime error from a dependency), a generic
// 500 INTERNAL_ERROR envelope is written that does NOT leak internals to the
// client — callers are expected to log the underlying error themselves.
//
// Returns the HTTP status code that was written, useful for logging and tests.
func FromError(w http.ResponseWriter, err error) int {
	if err == nil {
		// Defensive: shouldn't happen, but never write a success envelope on
		// the error path. Treat as 500.
		InternalError(w, "unknown error")
		return http.StatusInternalServerError
	}

	var de *phierr.DomainError
	if errors.As(err, &de) {
		status, ok := errorCodeToStatus[de.Code]
		if !ok {
			status = http.StatusInternalServerError
		}
		WriteJSON(w, status, Response{
			Success: false,
			Error: &Error{
				Code:    string(de.Code),
				Message: de.Message,
				Details: de.Details,
			},
		})
		return status
	}

	// Non-domain error: do not leak details to the caller. Caller should log.
	WriteJSON(w, http.StatusInternalServerError, Response{
		Success: false,
		Error: &Error{
			Code:    string(phierr.ErrInternal),
			Message: "internal server error",
		},
	})
	return http.StatusInternalServerError
}

// Unauthorized returns a 401 response.
func Unauthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = "unauthorized"
	}
	WriteJSON(w, http.StatusUnauthorized, Response{
		Success: false,
		Error: &Error{
			Code:    string(phierr.ErrUnauthorized),
			Message: message,
		},
	})
}

// Forbidden returns a 403 response.
func Forbidden(w http.ResponseWriter, message string) {
	if message == "" {
		message = "forbidden"
	}
	WriteJSON(w, http.StatusForbidden, Response{
		Success: false,
		Error: &Error{
			Code:    string(phierr.ErrForbidden),
			Message: message,
		},
	})
}

// Conflict returns a 409 response.
func Conflict(w http.ResponseWriter, message string) {
	if message == "" {
		message = "resource conflict"
	}
	WriteJSON(w, http.StatusConflict, Response{
		Success: false,
		Error: &Error{
			Code:    string(phierr.ErrConflict),
			Message: message,
		},
	})
}

// UnprocessableEntity returns a 422 response (semantic validation failure).
func UnprocessableEntity(w http.ResponseWriter, message string) {
	if message == "" {
		message = "unprocessable entity"
	}
	WriteJSON(w, http.StatusUnprocessableEntity, Response{
		Success: false,
		Error: &Error{
			Code:    string(phierr.ErrInvalidInput),
			Message: message,
		},
	})
}

package transport_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/philiaspace/phi-core/transport"
	phierr "github.com/philiaspace/phi-utils/errors"
)

func decodeEnvelope(t *testing.T, body []byte) transport.Response {
	t.Helper()
	var r transport.Response
	if err := json.Unmarshal(body, &r); err != nil {
		t.Fatalf("failed to decode envelope: %v\nbody=%s", err, string(body))
	}
	return r
}

func TestFromError_MapsDomainErrorCodes(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{"not_found", phierr.New(phierr.ErrNotFound, "missing"), http.StatusNotFound, "NOT_FOUND"},
		{"already_exists", phierr.New(phierr.ErrAlreadyExists, "dup"), http.StatusConflict, "ALREADY_EXISTS"},
		{"invalid_input", phierr.New(phierr.ErrInvalidInput, "bad"), http.StatusBadRequest, "INVALID_INPUT"},
		{"unauthorized", phierr.New(phierr.ErrUnauthorized, "no auth"), http.StatusUnauthorized, "UNAUTHORIZED"},
		{"forbidden", phierr.New(phierr.ErrForbidden, "denied"), http.StatusForbidden, "FORBIDDEN"},
		{"conflict", phierr.New(phierr.ErrConflict, "race"), http.StatusConflict, "CONFLICT"},
		{"internal", phierr.New(phierr.ErrInternal, "oops"), http.StatusInternalServerError, "INTERNAL_ERROR"},
		{"not_implemented", phierr.New(phierr.ErrNotImplemented, "soon"), http.StatusNotImplemented, "NOT_IMPLEMENTED"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			status := transport.FromError(rec, tc.err)
			if status != tc.wantStatus {
				t.Errorf("returned status: got %d want %d", status, tc.wantStatus)
			}
			if rec.Code != tc.wantStatus {
				t.Errorf("written status: got %d want %d", rec.Code, tc.wantStatus)
			}
			env := decodeEnvelope(t, rec.Body.Bytes())
			if env.Success {
				t.Error("expected Success=false on error path")
			}
			if env.Error == nil {
				t.Fatal("expected Error envelope, got nil")
			}
			if env.Error.Code != tc.wantCode {
				t.Errorf("error code: got %s want %s", env.Error.Code, tc.wantCode)
			}
		})
	}
}

func TestFromError_PreservesMessageAndDetails(t *testing.T) {
	err := phierr.New(phierr.ErrInvalidInput, "name is required").WithDetails("field=name")
	rec := httptest.NewRecorder()
	transport.FromError(rec, err)
	env := decodeEnvelope(t, rec.Body.Bytes())
	if env.Error.Message != "name is required" {
		t.Errorf("message: got %q", env.Error.Message)
	}
	if env.Error.Details != "field=name" {
		t.Errorf("details: got %q", env.Error.Details)
	}
}

func TestFromError_GenericErrorDoesNotLeakDetails(t *testing.T) {
	rec := httptest.NewRecorder()
	status := transport.FromError(rec, errors.New("sensitive internal detail: dsn=postgres://user:pass@host"))
	if status != http.StatusInternalServerError {
		t.Errorf("status: got %d want 500", status)
	}
	env := decodeEnvelope(t, rec.Body.Bytes())
	if env.Error.Message == "" {
		t.Fatal("expected a generic message")
	}
	// Must not contain the leaked detail
	if env.Error.Details != "" {
		t.Errorf("unexpected details on generic error: %q", env.Error.Details)
	}
	if env.Error.Code != "INTERNAL_ERROR" {
		t.Errorf("code: got %s want INTERNAL_ERROR", env.Error.Code)
	}
}

func TestFromError_NilErrorReturns500(t *testing.T) {
	rec := httptest.NewRecorder()
	status := transport.FromError(rec, nil)
	if status != http.StatusInternalServerError {
		t.Errorf("status: got %d want 500", status)
	}
}

func TestFromError_WrappedDomainError(t *testing.T) {
	wrapped := errors.Join(errors.New("context"), phierr.New(phierr.ErrNotFound, "missing"))
	rec := httptest.NewRecorder()
	status := transport.FromError(rec, wrapped)
	if status != http.StatusNotFound {
		t.Errorf("expected wrapped DomainError to map: got %d", status)
	}
	env := decodeEnvelope(t, rec.Body.Bytes())
	if env.Error.Code != "NOT_FOUND" {
		t.Errorf("code: got %s", env.Error.Code)
	}
}

func TestHelpers_StatusCodes(t *testing.T) {
	cases := []struct {
		name   string
		fn     func(http.ResponseWriter, string)
		want   int
		code   string
	}{
		{"unauthorized", transport.Unauthorized, http.StatusUnauthorized, "UNAUTHORIZED"},
		{"forbidden", transport.Forbidden, http.StatusForbidden, "FORBIDDEN"},
		{"conflict", transport.Conflict, http.StatusConflict, "CONFLICT"},
		{"unprocessable", transport.UnprocessableEntity, http.StatusUnprocessableEntity, "INVALID_INPUT"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tc.fn(rec, "msg")
			if rec.Code != tc.want {
				t.Errorf("status: got %d want %d", rec.Code, tc.want)
			}
			env := decodeEnvelope(t, rec.Body.Bytes())
			if env.Error == nil || env.Error.Code != tc.code {
				t.Errorf("code: got %+v want %s", env.Error, tc.code)
			}
			if env.Success {
				t.Error("expected Success=false")
			}
		})
	}
}

func TestHelpers_DefaultMessageWhenEmpty(t *testing.T) {
	rec := httptest.NewRecorder()
	transport.Forbidden(rec, "")
	env := decodeEnvelope(t, rec.Body.Bytes())
	if env.Error.Message != "forbidden" {
		t.Errorf("expected default message, got %q", env.Error.Message)
	}
}

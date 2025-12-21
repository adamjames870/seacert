package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamjames870/seacert/internal/api/handlers"
)

func TestRespondWithJSON(t *testing.T) {
	w := httptest.NewRecorder()
	payload := map[string]string{"message": "hello"}
	code := http.StatusOK

	err := handlers.RespondWithJSON(w, code, payload)
	if err != nil {
		t.Errorf("RespondWithJSON returned error: %v", err)
	}

	if w.Code != code {
		t.Errorf("Expected status code %d, got %d", code, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var got map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if got["message"] != "hello" {
		t.Errorf("Expected message 'hello', got '%s'", got["message"])
	}
}

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name string
		code int
		msg  string
		err  error
	}{
		{"Client Error", http.StatusBadRequest, "bad request", errors.New("invalid input")},
		{"Server Error", http.StatusInternalServerError, "internal error", errors.New("db failure")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := handlers.RespondWithError(w, tt.code, tt.msg, tt.err)
			if err != nil {
				t.Errorf("RespondWithError returned error: %v", err)
			}

			if w.Code != tt.code {
				t.Errorf("Expected status code %d, got %d", tt.code, w.Code)
			}

			var got map[string]string
			json.Unmarshal(w.Body.Bytes(), &got)
			if got["error"] != tt.msg {
				t.Errorf("Expected error message '%s', got '%s'", tt.msg, got["error"])
			}
		})
	}
}

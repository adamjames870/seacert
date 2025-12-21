package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResponseWriter(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec, status: http.StatusOK}

	rw.WriteHeader(http.StatusCreated)
	if rw.status != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rw.status)
	}

	// Test that double WriteHeader is ignored
	rw.WriteHeader(http.StatusOK)
	if rw.status != http.StatusCreated {
		t.Errorf("expected status to remain %d, got %d", http.StatusCreated, rw.status)
	}

	rw.Write([]byte("hello"))
	if rec.Body.String() != "hello" {
		t.Errorf("expected body 'hello', got %s", rec.Body.String())
	}
}

func TestLoggingMiddleware(t *testing.T) {
	t.Run("Empty body", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("ok"))
		})

		loggingHandler := Logging(handler)

		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()

		loggingHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusAccepted {
			t.Errorf("expected status %d, got %d", http.StatusAccepted, rec.Code)
		}
	})

	t.Run("With body", func(t *testing.T) {
		bodyText := "test body content"
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify body is still readable
			body, _ := io.ReadAll(r.Body)
			if string(body) != bodyText {
				t.Errorf("expected body %q, got %q", bodyText, string(body))
			}
			w.WriteHeader(http.StatusOK)
		})

		loggingHandler := Logging(handler)

		req := httptest.NewRequest("POST", "/test", strings.NewReader(bodyText))
		rec := httptest.NewRecorder()

		loggingHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})
}

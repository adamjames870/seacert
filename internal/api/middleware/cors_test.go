package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCors(t *testing.T) {
	tests := []struct {
		name           string
		allowedOrigins string
		originHeader   string
		method         string
		expectedOrigin string
		expectedStatus int
	}{
		{
			name:           "Default allow all when empty",
			allowedOrigins: "",
			originHeader:   "http://example.com",
			method:         "GET",
			expectedOrigin: "*",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Match specific origin",
			allowedOrigins: "http://localhost:3000,http://example.com",
			originHeader:   "http://example.com",
			method:         "GET",
			expectedOrigin: "http://example.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "No match",
			allowedOrigins: "http://localhost:3000",
			originHeader:   "http://example.com",
			method:         "GET",
			expectedOrigin: "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "OPTIONS request",
			allowedOrigins: "",
			originHeader:   "http://example.com",
			method:         "OPTIONS",
			expectedOrigin: "*",
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ALLOWED_ORIGINS", tt.allowedOrigins)
			defer os.Unsetenv("ALLOWED_ORIGINS")

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handler := Cors(nextHandler)
			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.originHeader != "" {
				req.Header.Set("Origin", tt.originHeader)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			origin := rec.Header().Get("Access-Control-Allow-Origin")
			if origin != tt.expectedOrigin {
				t.Errorf("expected origin %q, got %q", tt.expectedOrigin, origin)
			}
		})
	}
}

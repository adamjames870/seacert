package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/adamjames870/seacert/internal/api/auth"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var body []byte
		if r.Body != nil {
			var err error
			// Limit the amount we read for logging to 2KB
			const maxLogSize = 2048
			body, err = io.ReadAll(io.LimitReader(r.Body, maxLogSize))
			if err != nil {
				slog.Error("Failed to read request body", "error", err)
			}
			r.Body = io.NopCloser(io.MultiReader(bytes.NewBuffer(body), r.Body))
		}

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		logger := slog.Default()
		if user, ok := auth.UserFromContext(r.Context()); ok {
			logger = logger.With("user_id", user.Id, "email", user.Email)
		}

		logger.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration", duration,
			"ip", r.RemoteAddr,
		)

		if len(body) > 0 {
			logger.Debug("HTTP request body", "body", string(body))
		}
	})
}

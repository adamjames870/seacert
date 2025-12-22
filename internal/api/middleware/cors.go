package middleware

import (
	"net/http"
	"os"
	"strings"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		origin := r.Header.Get("Origin")

		if allowedOrigins == "" {
			// Default to allow all if not specified (or you might want to be stricter)
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			origins := strings.Split(allowedOrigins, ",")
			for _, o := range origins {
				if o == origin || o == "*" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

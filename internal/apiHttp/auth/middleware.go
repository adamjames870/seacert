package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoAuthHeader = errors.New("no authorization header provided")
	ErrInvalidToken = errors.New("invalid token")
)

func Middleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, ErrNoAuthHeader.Error(), http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Parse and validate token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Supabase access tokens are HS256 signed
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, ErrInvalidToken
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			// Extract claims (user info)
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			// Optional: check token expiration
			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					http.Error(w, "token expired", http.StatusUnauthorized)
					return
				}
			}

			maxTokenAge := int64(3600) // 1 hour
			if iat, ok := claims["iat"].(float64); ok {
				if time.Now().Unix()-int64(iat) > maxTokenAge {
					http.Error(w, "token too old", http.StatusUnauthorized)
					return
				}
			}

			if aud, ok := claims["aud"].(string); !ok || aud != "authenticated" {
				http.Error(w, "invalid audience", http.StatusUnauthorized)
				return
			}

			if role, ok := claims["role"].(string); !ok || role != "authenticated" {
				http.Error(w, "invalid role", http.StatusUnauthorized)
				return
			}

			user := User{
				ID:    claims["sub"].(string),
				Email: claims["email"].(string),
				Role:  claims["role"].(string),
			}

			// Store claims in context for handlers
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

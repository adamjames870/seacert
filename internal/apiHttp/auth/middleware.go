package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	ErrNoAuthHeader = errors.New("no authorization header provided")
	ErrInvalidToken = errors.New("invalid token")
)

type ctxKey string

const UserIDKey ctxKey = "userID"

func loadSupabaseJWK(apiKey string) (jwk.Key, error) {

	key, err := jwk.ParseKey([]byte(apiKey))
	if err != nil {
		return nil, err
	}

	return key, nil
}

func Middleware(authInfo Info) (func(http.Handler) http.Handler, error) {

	key, err := loadSupabaseJWK(authInfo.ApiKey)
	if err != nil {
		return nil, err
	}

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
			token, err := jwt.ParseString(
				tokenString,
				jwt.WithKey(jwa.ES256, key),
				jwt.WithValidate(true),
				jwt.WithIssuer(authInfo.ExpectedIssuer),
				jwt.WithAudience(authInfo.ExpectedAudience),
			)
			if err != nil {
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			issuedAt := token.IssuedAt()
			now := time.Now()
			if now.Sub(issuedAt) > time.Hour {
				http.Error(w, "token too old", http.StatusUnauthorized)
				return
			}

			role, ok := token.Get("role")
			if !ok || role.(string) != "authenticated" {
				http.Error(w, "invalid role", http.StatusUnauthorized)
				return
			}

			email, ok := token.Get("email")
			if !ok || email == nil {
				http.Error(w, "invalid email", http.StatusUnauthorized)
				return
			}

			user := User{
				ID:    token.Subject(),
				Email: email.(string),
				Role:  role.(string),
			}

			// Store claims in context for handlers
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}, nil
}

package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/domain/users"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	ErrNoAuthHeader = errors.New("no authorization header provided")
	ErrInvalidToken = errors.New("invalid token")
)

func NewAuthMiddleware(authInfo Info, state *internal.ApiState) (func(http.Handler) http.Handler, error) {

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

			const bearerPrefix = "Bearer "

			if !strings.HasPrefix(authHeader, bearerPrefix) {
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
			if tokenString == "" {
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

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

			user, errUser := userFromToken(token)
			if errUser != nil {
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			uuidId, errParse := uuid.Parse(user.Id)
			if errParse != nil {
				http.Error(w, "user id is not a valid uuid", http.StatusBadRequest)
				return
			}

			_, errUserExists := users.EnsureUserExists(state, r.Context(), uuidId, user.Email)
			if errUserExists != nil {
				http.Error(w, "user cannot be found or created", http.StatusBadRequest)
			}

			// Store claims in context for handlers
			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}, nil
}

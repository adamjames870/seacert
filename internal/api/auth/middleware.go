package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/adamjames870/seacert/internal/domain/users"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	ErrNoAuthHeader   = errors.New("no authorization header provided")
	ErrInvalidToken   = errors.New("invalid token")
	ErrInternalServer = errors.New("internal server error")
)

type UserProvider interface {
	EnsureUserExists(ctx context.Context, id uuid.UUID, email string) (users.User, error)
}

func NewAuthMiddleware(authInfo Info, userStore UserProvider) (func(http.Handler) http.Handler, error) {

	key, err := loadSupabaseJWK(authInfo.PublicKey)
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

			token, err := jwt.ParseString(
				tokenString,
				jwt.WithKey(key.Algorithm(), key),
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
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			_, errUserExists := userStore.EnsureUserExists(r.Context(), uuidId, user.Email)
			if errUserExists != nil {
				http.Error(w, ErrInternalServer.Error(), http.StatusInternalServerError)
				return
			}

			// Store claims in context for handlers
			ctx := WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}, nil
}

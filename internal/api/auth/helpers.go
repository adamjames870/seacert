package auth

import (
	"context"
	"errors"

	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	ErrInvalidRole       = errors.New("invalid role")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrMissingSubject    = errors.New("missing subject")
	ErrUserNotFoundInCtx = errors.New("user not found in context")
)

func getStringClaim(t jwt.Token, name string) (string, bool) {
	v, ok := t.Get(name)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok && s != ""
}

func userFromToken(t jwt.Token) (dto.User, error) {
	role, ok := getStringClaim(t, "role")
	if !ok || role != "authenticated" {
		return dto.User{}, ErrInvalidRole
	}

	email, ok := getStringClaim(t, "email")
	if !ok {
		return dto.User{}, ErrInvalidEmail
	}

	sub := t.Subject()
	if sub == "" {
		return dto.User{}, ErrMissingSubject
	}

	id, errId := uuid.Parse(sub)
	if errId != nil {
		return dto.User{}, errId
	}

	// Extract user role from app_metadata
	userRole := "user" // default role
	appMetadata, ok := t.Get("app_metadata")
	if ok {
		if metadata, ok := appMetadata.(map[string]interface{}); ok {
			if r, ok := metadata["role"].(string); ok {
				userRole = r
			}
		}
	}

	return dto.User{
		Id:    id.String(),
		Email: email,
		Role:  userRole,
	}, nil
}

func loadSupabaseJWK(jwkStr string) (jwk.Key, error) {
	key, err := jwk.ParseKey([]byte(jwkStr))
	if err != nil {
		return nil, err
	}

	// Ensure the algorithm is set, default to ES256 if not specified
	if key.Algorithm().String() == "" {
		key.Set(jwk.AlgorithmKey, "ES256")
	}

	return key, nil
}

type contextKey struct{}

var userKey = contextKey{}

func WithUser(ctx context.Context, user dto.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) (dto.User, bool) {
	user, ok := ctx.Value(userKey).(dto.User)
	return user, ok
}

func UserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	authUser, ok := UserFromContext(ctx)
	if !ok {
		return uuid.UUID{}, ErrUserNotFoundInCtx
	}

	uuidId, errParse := uuid.Parse(authUser.Id)
	if errParse != nil {
		return uuid.UUID{}, errParse
	}

	return uuidId, nil
}

package auth

import (
	"context"
	"errors"

	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func loadSupabaseJWK(apiKey string) (jwk.Key, error) {

	key, err := jwk.ParseKey([]byte(apiKey))
	if err != nil {
		return nil, err
	}

	return key, nil
}

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
		return dto.User{}, errors.New("invalid role")
	}

	email, ok := getStringClaim(t, "email")
	if !ok {
		return dto.User{}, errors.New("invalid email")
	}

	sub := t.Subject()
	if sub == "" {
		return dto.User{}, errors.New("missing subject")
	}

	id, errId := uuid.Parse(sub)
	if errId != nil {
		return dto.User{}, errId
	}

	return dto.User{
		Id:    id.String(),
		Email: email,
	}, nil
}

type userContextKeyType string

const userContextKey userContextKeyType = "user"

func UserFromContext(ctx context.Context) (dto.User, bool) {
	user, ok := ctx.Value(userContextKey).(dto.User)
	return user, ok
}

func UserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	authUser, ok := UserFromContext(ctx)
	if !ok {
		return uuid.UUID{}, errors.New("user not found in context")
	}

	uuidId, errParse := uuid.Parse(authUser.Id)
	if errParse != nil {
		return uuid.UUID{}, errParse
	}

	return uuidId, nil
}

package auth

import (
	"context"
	"errors"

	"github.com/adamjames870/seacert/internal/domain/users"
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

func userFromToken(t jwt.Token) (users.User, error) {
	role, ok := getStringClaim(t, "role")
	if !ok || role != "authenticated" {
		return users.User{}, errors.New("invalid role")
	}

	email, ok := getStringClaim(t, "email")
	if !ok {
		return users.User{}, errors.New("invalid email")
	}

	sub := t.Subject()
	if sub == "" {
		return users.User{}, errors.New("missing subject")
	}

	id, errId := uuid.Parse(sub)
	if errId != nil {
		return users.User{}, errId
	}

	return users.User{
		Id:    id,
		Email: email,
	}, nil
}

func UserFromContext(ctx context.Context) (users.User, bool) {
	user, ok := ctx.Value(userContextKey).(users.User)
	return user, ok
}

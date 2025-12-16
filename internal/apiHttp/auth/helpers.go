package auth

import (
	"errors"

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

func userFromToken(t jwt.Token) (User, error) {
	role, ok := getStringClaim(t, "role")
	if !ok || role != "authenticated" {
		return User{}, errors.New("invalid role")
	}

	email, ok := getStringClaim(t, "email")
	if !ok {
		return User{}, errors.New("invalid email")
	}

	sub := t.Subject()
	if sub == "" {
		return User{}, errors.New("missing subject")
	}

	return User{
		ID:    sub,
		Email: email,
		Role:  role,
	}, nil
}

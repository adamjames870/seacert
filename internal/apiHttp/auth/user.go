package auth

import "context"

type userContextKeyType string

const userContextKey userContextKeyType = "user"

type User struct {
	ID    string
	Email string
	Role  string
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey).(User)
	return user, ok
}

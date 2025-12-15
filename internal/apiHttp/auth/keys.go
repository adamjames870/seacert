package auth

type contextKey string

const UserContextKey contextKey = "user"

type User struct {
	ID    string
	Email string
	Role  string
}

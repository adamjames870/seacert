package auth

type contextKey string

const JwksUrl = "https://tszamyytdpalngxtzmpt.supabase.co/auth/v1/jwks"

const UserContextKey contextKey = "user"

type User struct {
	ID    string
	Email string
	Role  string
}

type Info struct {
	JwksUrl          string
	ApiKey           string
	ExpectedIssuer   string
	ExpectedAudience string
}

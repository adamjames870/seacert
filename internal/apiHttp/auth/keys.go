package auth

type userContextKeyType string

const userContextKey userContextKeyType = "user"

type Info struct {
	JwksUrl          string
	ApiKey           string
	ExpectedIssuer   string
	ExpectedAudience string
}

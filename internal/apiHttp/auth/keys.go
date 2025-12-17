package auth

type Info struct {
	JwksUrl          string
	ApiKey           string
	ExpectedIssuer   string
	ExpectedAudience string
}

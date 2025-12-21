package api

import (
	"net/http"
	"os"
	"testing"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/logging"
)

func TestCreateEndpoints(t *testing.T) {
	// Set required environment variables for NewAuthMiddleware which is called in createEndpoints
	os.Setenv("SUPABASE_PUBLIC_JWK", `{"kty":"RSA","n":"0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93laj7XJhP2MAnwy7UZnnWQS9geS4_302z1nw5o8tC9KxVOTB4T8r412EsqI5v44YFvK4L4yM1h-444444444444444444444444444444444444444444444444444444444444444444444444444444444","e":"AQAB"}`)
	defer os.Unsetenv("SUPABASE_PUBLIC_JWK")

	mux := http.NewServeMux()
	state := &internal.ApiState{
		Logger: logging.NewLogger(),
	}

	err := createEndpoints(mux, state)
	if err != nil {
		t.Fatalf("createEndpoints failed: %v", err)
	}

	// We can't easily check if specific routes were added to http.ServeMux in older Go versions,
	// but in Go 1.22+ we might be able to. The project uses go 1.25.4 (based on go.mod).
	// However, the main point is to ensure it doesn't panic or return an error.
}

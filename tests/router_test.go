package tests

import (
	"os"
	"testing"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api"
	"github.com/adamjames870/seacert/internal/logging"
)

func TestBuildRouter(t *testing.T) {
	// Set required environment variables for BuildRouter
	os.Setenv("SUPABASE_PUBLIC_JWK", "{\"kty\":\"EC\",\"crv\":\"P-256\",\"x\":\"MKBfRE2_8W3_Z_pD0v3_F_Z_pD0v3_F_Z_pD0v3_F_0\",\"y\":\"MKBfRE2_8W3_Z_pD0v3_F_Z_pD0v3_F_Z_pD0v3_F_0\"}")
	os.Setenv("SUPABASE_ISSUER", "https://example.supabase.co/auth/v1")
	os.Setenv("SUPABASE_AUDIENCE", "authenticated")
	defer os.Unsetenv("SUPABASE_PUBLIC_JWK")
	defer os.Unsetenv("SUPABASE_ISSUER")
	defer os.Unsetenv("SUPABASE_AUDIENCE")

	state := &internal.ApiState{
		Logger: logging.NewLogger(),
		IsDev:  true,
	}

	handler, err := api.BuildRouter(state)
	if err != nil {
		t.Fatalf("BuildRouter failed: %v", err)
	}

	if handler == nil {
		t.Fatal("BuildRouter returned nil handler")
	}
}

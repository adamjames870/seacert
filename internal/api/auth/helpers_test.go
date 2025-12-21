package auth

import (
	"testing"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

func TestGetStringClaim(t *testing.T) {
	tok := jwt.New()
	tok.Set("foo", "bar")
	tok.Set("empty", "")
	tok.Set("notstring", 123)

	tests := []struct {
		name string
		key  string
		want string
		ok   bool
	}{
		{"valid claim", "foo", "bar", true},
		{"empty claim", "empty", "", false},
		{"missing claim", "missing", "", false},
		{"not a string", "notstring", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := getStringClaim(tok, tt.key)
			if got != tt.want || ok != tt.ok {
				t.Errorf("getStringClaim() = (%v, %v), want (%v, %v)", got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestLoadSupabaseJWK(t *testing.T) {
	// A simple valid RSA public key in JWK format for testing
	validJWK := `{"kty":"RSA","n":"0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93laj7XJhP2MAnwy7UZnnWQS9geS4_302z1nw5o8tC9KxVOTB4T8r412EsqI5v44YFvK4L4yM1h-444444444444444444444444444444444444444444444444444444444444444444444444444444444","e":"AQAB"}`

	t.Run("valid JWK", func(t *testing.T) {
		key, err := loadSupabaseJWK(validJWK)
		if err != nil {
			t.Fatalf("loadSupabaseJWK failed: %v", err)
		}
		if key == nil {
			t.Fatal("key is nil")
		}
		// Check if algorithm is set to ES256 by default if not present
		if key.Algorithm().String() != "ES256" {
			t.Errorf("expected algorithm ES256, got %s", key.Algorithm().String())
		}
	})

	t.Run("invalid JWK", func(t *testing.T) {
		_, err := loadSupabaseJWK("invalid")
		if err == nil {
			t.Error("expected error for invalid JWK, got nil")
		}
	})
}

func TestUserFromToken(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name    string
		setup   func() jwt.Token
		wantErr error
	}{
		{
			name: "valid token",
			setup: func() jwt.Token {
				tok := jwt.New()
				tok.Set("role", "authenticated")
				tok.Set("email", "test@example.com")
				tok.Set("sub", validUUID)
				return tok
			},
			wantErr: nil,
		},
		{
			name: "invalid role",
			setup: func() jwt.Token {
				tok := jwt.New()
				tok.Set("role", "other")
				return tok
			},
			wantErr: ErrInvalidRole,
		},
		{
			name: "missing email",
			setup: func() jwt.Token {
				tok := jwt.New()
				tok.Set("role", "authenticated")
				tok.Set("sub", validUUID)
				return tok
			},
			wantErr: ErrInvalidEmail,
		},
		{
			name: "missing subject",
			setup: func() jwt.Token {
				tok := jwt.New()
				tok.Set("role", "authenticated")
				tok.Set("email", "test@example.com")
				return tok
			},
			wantErr: ErrMissingSubject,
		},
		{
			name: "invalid uuid subject",
			setup: func() jwt.Token {
				tok := jwt.New()
				tok.Set("role", "authenticated")
				tok.Set("email", "test@example.com")
				tok.Set("sub", "not-a-uuid")
				return tok
			},
			wantErr: nil, // uuid.Parse will return error, but not ErrMissingSubject
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := tt.setup()
			user, err := userFromToken(tok)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("userFromToken() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if tt.name == "invalid uuid subject" {
					if err == nil {
						t.Error("expected error for invalid UUID, got nil")
					}
				} else if err != nil {
					t.Errorf("userFromToken() unexpected error = %v", err)
				} else {
					if user.Email != "test@example.com" {
						t.Errorf("expected email test@example.com, got %s", user.Email)
					}
					if user.Id != validUUID {
						t.Errorf("expected id %s, got %s", validUUID, user.Id)
					}
				}
			}
		})
	}
}

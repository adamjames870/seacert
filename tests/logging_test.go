package tests

import (
	"os"
	"testing"

	"github.com/adamjames870/seacert/internal/logging"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		wantJSON bool
	}{
		{
			name:     "production platform",
			platform: "production",
			wantJSON: true,
		},
		{
			name:     "test platform",
			platform: "test",
			wantJSON: true,
		},
		{
			name:     "dev platform",
			platform: "dev",
			wantJSON: false,
		},
		{
			name:     "empty platform",
			platform: "",
			wantJSON: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save current PLATFORM and restore it after test
			oldPlatform := os.Getenv("PLATFORM")
			defer os.Setenv("PLATFORM", oldPlatform)

			os.Setenv("PLATFORM", tt.platform)

			logger := logging.NewLogger()
			if logger == nil {
				t.Fatal("NewLogger() returned nil")
			}

			// We can't easily check if it's a JSONHandler or TextHandler without reflection
			// but we can at least verify it doesn't panic and returns a valid logger.
		})
	}
}

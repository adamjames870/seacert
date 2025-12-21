package main

import (
	"os"
	"testing"

	"github.com/adamjames870/seacert/internal"
)

func TestSetDevFlag(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		want     bool
	}{
		{
			name:     "platform is dev",
			platform: "dev",
			want:     true,
		},
		{
			name:     "platform is production",
			platform: "production",
			want:     false,
		},
		{
			name:     "platform is empty",
			platform: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldPlatform := os.Getenv("PLATFORM")
			defer os.Setenv("PLATFORM", oldPlatform)

			os.Setenv("PLATFORM", tt.platform)

			state := &internal.ApiState{}
			setDevFlag(state)

			if state.IsDev != tt.want {
				t.Errorf("setDevFlag() state.IsDev = %v, want %v", state.IsDev, tt.want)
			}
		})
	}
}

func TestLoadDb_NoUrl(t *testing.T) {
	oldUrl := os.Getenv("DB_URL")
	defer os.Setenv("DB_URL", oldUrl)

	os.Unsetenv("DB_URL")

	state := &internal.ApiState{}
	err := loadDb(state)

	if err == nil {
		t.Fatal("expected error when DB_URL is not set, got nil")
	}

	expectedErr := "DB_URL environment variable is not set"
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestLoadState_Fail(t *testing.T) {
	oldUrl := os.Getenv("DB_URL")
	defer os.Setenv("DB_URL", oldUrl)
	os.Unsetenv("DB_URL")

	state := &internal.ApiState{}
	err := LoadState(state)

	if err == nil {
		t.Fatal("expected error from LoadState when DB_URL is missing, got nil")
	}
}

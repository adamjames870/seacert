package main

import (
	"os"
	"testing"

	"github.com/adamjames870/seacert/internal"
)

func TestRun_Shutdown(t *testing.T) {
	// Skip if no DB_URL is available as run calls LoadState which calls loadDb
	// and loadDb pings the database.
	// Actually, we can try to test it by providing a mock/fake DB URL if the driver supports it,
	// but here it's "postgres".

	// If we want to test 'run' successfully, we need a working DB or it will fail at LoadState.
	// Let's test the failure case when DB is missing.

	oldUrl := os.Getenv("DB_URL")
	defer os.Setenv("DB_URL", oldUrl)
	os.Unsetenv("DB_URL")

	state := &internal.ApiState{}
	err := run(state)

	if err == nil {
		t.Fatal("expected error from run when DB_URL is missing, got nil")
	}
}

func TestRun_Signal(t *testing.T) {
	// This test would be complex because 'run' blocks until a signal is received.
	// We could run it in a goroutine and send a signal, but we need a valid state first.
	// Since we can't easily mock the DB ping in loadDb without changing code,
	// we'll stick to testing what we can.
}

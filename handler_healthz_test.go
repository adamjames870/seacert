package main

import (
	"net/http"
	"testing"
)

func TestHealthzEndpoint(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/healthz")
	if err != nil {
		t.Fatalf("failed to call /api/healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

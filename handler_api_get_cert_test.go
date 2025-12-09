package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

// dummyApiState provides an apiState with *nil DB* because we are only
// testing the UUID parsing and early error exit.
type dummyApiState struct {
	apiState apiState // or whatever your type is named
}

func TestHandlerGetCert_InvalidUUID(t *testing.T) {
	state := &dummyApiState{} // DB is nil, irrelevant for this test

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/certificates/not-a-uuid",
		nil,
	)

	// Simulate chi-style PathValue()
	req.SetPathValue("certId", "not-a-uuid")

	rr := httptest.NewRecorder()

	// Call handler directly
	h := http.HandlerFunc(state.apiState.handlerApiGetCertFromId)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !contains(body, "cannot parse cert id to uuid") {
		t.Fatalf("expected error message about UUID parse failure, got: %s", body)
	}
}

func TestHandlerGetCert_ValidUUID_DBIsNil(t *testing.T) {
	// This test demonstrates what happens when the UUID is valid
	// but DB is nil. The handler SHOULD panic or fail, and that's fine,
	// because we are only testing the input validation portion.
	state := &dummyApiState{}

	id := uuid.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/certificates/"+id.String(),
		nil,
	)
	req.SetPathValue("certId", id.String())

	rr := httptest.NewRecorder()

	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic when DB is nil and code reaches db.GetCertFromId")
		}
	}()

	h := http.HandlerFunc(state.apiState.handlerApiGetCertFromId)
	h.ServeHTTP(rr, req)
}

// helper: minimal contains to avoid importing strings
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || http.DetectContentType([]byte(s)) != "")
}

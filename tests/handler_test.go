package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamjames870/seacert/internal/api/handlers/admin"
)

func TestHandlerAdminHealthz(t *testing.T) {
	handler := admin.HandlerAdminHealthz()

	req, err := http.NewRequest("GET", "/admin/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "OK\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	contentType := rr.Header().Get("Content-Type")
	expectedContentType := "text/plain; charset=utf-8"
	if contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}
}

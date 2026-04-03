package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/dto"
)

func TestDecodeAndValidate_ParamsAddCertificate(t *testing.T) {
	// UserId should NOT be required in the JSON payload as it's added by the server
	payload := map[string]interface{}{
		"cert-type-id": "some-id",
		"cert-number":  "12345",
		"issuer-id":    "issuer-id",
		"issued-date":  "2023-01-01",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	var params dto.ParamsAddCertificate
	err := handlers.DecodeAndValidate(req, &params)

	if err != nil {
		t.Errorf("DecodeAndValidate failed: %v", err)
	}

	if params.CertTypeId != "some-id" {
		t.Errorf("Expected CertTypeId 'some-id', got '%s'", params.CertTypeId)
	}
}

func TestDecodeAndValidate_ParamsAddCertificate_MissingFields(t *testing.T) {
	// Missing cert-number, which IS required
	payload := map[string]interface{}{
		"cert-type-id": "some-id",
		"issuer-id":    "issuer-id",
		"issued-date":  "2023-01-01",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	var params dto.ParamsAddCertificate
	err := handlers.DecodeAndValidate(req, &params)

	if err == nil {
		t.Error("Expected validation error for missing cert-number, got nil")
	}
}

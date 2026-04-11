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

func TestDecodeAndValidate_ParamsAddCertificateType_ZeroValidity(t *testing.T) {
	payload := map[string]any{
		"name":                   "Test Cert",
		"short-name":             "TC",
		"normal-validity-months": 0,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/cert-types", bytes.NewReader(body))

	var params dto.ParamsAddCertificateType
	err := handlers.DecodeAndValidate(req, &params)

	if err != nil {
		t.Errorf("DecodeAndValidate failed for zero validity: %v", err)
	}

	if params.NormalValidityMonths == nil {
		t.Fatal("Expected NormalValidityMonths to be non-nil")
	}

	if *params.NormalValidityMonths != 0 {
		t.Errorf("Expected NormalValidityMonths to be 0, got %d", *params.NormalValidityMonths)
	}
}

func TestDecodeAndValidate_ParamsAddCertificateType_MissingValidity(t *testing.T) {
	payload := map[string]any{
		"name":       "Test Cert",
		"short-name": "TC",
		// missing normal-validity-months
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/cert-types", bytes.NewReader(body))

	var params dto.ParamsAddCertificateType
	err := handlers.DecodeAndValidate(req, &params)

	if err == nil {
		t.Error("Expected validation error for missing normal-validity-months, got nil")
	}
}

package certificates

import (
	"encoding/json"
	"testing"

	"github.com/adamjames870/seacert/internal/dto"
)

func TestUnmarshalExtractedCertificate(t *testing.T) {
	jsonData := `{
		"cert-type-name": "Master Mariner",
		"cert-number": "123456",
		"issuer-name": "UK MCA",
		"issued-date": "2023-01-01",
		"expiry-date": "2028-01-01",
		"remarks": "Valid for all ships",
		"cert-type-id": "cert-uuid-123",
		"issuer-id": "issuer-uuid-456"
	}`

	var extracted dto.ExtractedCertificate
	err := json.Unmarshal([]byte(jsonData), &extracted)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if extracted.CertTypeName != "Master Mariner" {
		t.Errorf("Expected CertTypeName 'Master Mariner', got '%s'", extracted.CertTypeName)
	}
	if extracted.CertNumber != "123456" {
		t.Errorf("Expected CertNumber '123456', got '%s'", extracted.CertNumber)
	}
	if extracted.IssuerName != "UK MCA" {
		t.Errorf("Expected IssuerName 'UK MCA', got '%s'", extracted.IssuerName)
	}
	if extracted.IssuedDate != "2023-01-01" {
		t.Errorf("Expected IssuedDate '2023-01-01', got '%s'", extracted.IssuedDate)
	}
	if extracted.ExpiryDate == nil || *extracted.ExpiryDate != "2028-01-01" {
		t.Errorf("Expected ExpiryDate '2028-01-01', got '%v'", extracted.ExpiryDate)
	}
	if extracted.CertTypeId == nil || *extracted.CertTypeId != "cert-uuid-123" {
		t.Errorf("Expected CertTypeId 'cert-uuid-123', got '%v'", extracted.CertTypeId)
	}
	if extracted.IssuerId == nil || *extracted.IssuerId != "issuer-uuid-456" {
		t.Errorf("Expected IssuerId 'issuer-uuid-456', got '%v'", extracted.IssuerId)
	}
}

func TestUnmarshalExtractedCertificateWithNulls(t *testing.T) {
	jsonData := `{
		"cert-type-name": "Master Mariner",
		"cert-number": "123456",
		"issuer-name": "UK MCA",
		"issued-date": "2023-01-01",
		"expiry-date": null,
		"remarks": null
	}`

	var extracted dto.ExtractedCertificate
	err := json.Unmarshal([]byte(jsonData), &extracted)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if extracted.ExpiryDate != nil {
		t.Errorf("Expected ExpiryDate to be nil, got '%v'", *extracted.ExpiryDate)
	}
	if extracted.Remarks != nil {
		t.Errorf("Expected Remarks to be nil, got '%v'", *extracted.Remarks)
	}
}

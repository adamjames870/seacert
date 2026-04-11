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

func TestUnmarshalExtractedCertificateArray(t *testing.T) {
	jsonData := `[
	  {
		"cert-type-name": "STCW UPDATED PROFICIENCY IN FIRE PREVENTION & FIRE FIGHTING",
		"cert-number": "MSA-15886",
		"issuer-name": "Maritime Skills Academy (Dover)",
		"issued-date": "2025-02-18",
		"expiry-date": null,
		"remarks": "STCW Reg. VI/1 (para 1) Sec. A-VI/1 (para 3 and 4.2). Date of Birth: 26-Aug-1983.",
		"cert-type-id": "2e81e015-99cc-478e-bea8-3597925708b6",
		"issuer-id": "56727bea-f002-4492-882d-c9084f9229e0"
	  }
	]`

	// Since we can't easily call ExtractCertificateData without a real genai client,
	// we'll test the logic that was added to it here.
	var extracted dto.ExtractedCertificate
	rawText := []byte(jsonData)

	// This mimics the logic in ExtractCertificateData
	if err := json.Unmarshal(rawText, &extracted); err != nil {
		var list []dto.ExtractedCertificate
		if errArray := json.Unmarshal(rawText, &list); errArray == nil && len(list) > 0 {
			extracted = list[0]
		} else {
			t.Fatalf("Failed to unmarshal as single object or array: %v", err)
		}
	}

	if extracted.CertNumber != "MSA-15886" {
		t.Errorf("Expected CertNumber 'MSA-15886', got '%s'", extracted.CertNumber)
	}
	if extracted.CertTypeId == nil || *extracted.CertTypeId != "2e81e015-99cc-478e-bea8-3597925708b6" {
		t.Errorf("Expected CertTypeId '2e81e015-99cc-478e-bea8-3597925708b6', got '%v'", extracted.CertTypeId)
	}
}

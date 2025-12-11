package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/adamjames870/seacert/models"
	"github.com/google/uuid"
)

var createdCertType models.CertificateType
var createdCert models.Certificate

func StoreCreatedCertificateType(certType models.CertificateType) {
	createdCertType = certType
}

func StoreCreatedCertificate(cert models.Certificate) {
	createdCert = cert
}

func TestInOrder(t *testing.T) {

	fmt.Println("Running integration tests...")

	tHealthzEndpoint(t)
	tResetDb(t)
	//tHandlerAddCertType(t)
	//tHandlerAddCertificate(t)
	//tGetCertificateFromId(t)
	//
	//fmt.Println("Initial integration tests complete. Installing dummy certs...")
	//
	//err := LoadDummyCerts()
	//if err != nil {
	//	fmt.Println("Failed to load dummy certs:" + err.Error())
	//}
	//
	//fmt.Println("Dummy certs installed. Checking retrieval...")
	//
	//tGetAllCerts(t)
	//
	//fmt.Println("All integration tests run.")

}

func tHealthzEndpoint(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/admin/healthz")
	if err != nil {
		t.Fatalf("failed to call /api/healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("HEALTHZ: expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func tResetDb(t *testing.T) {
	resp, err := http.Post("http://localhost:8080/admin/reset", "", nil)
	if err != nil {

		t.Fatalf("failed to call /api/reset: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("RESET: expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func tHandlerAddCertType(t *testing.T) {

	url := "http://localhost:8080/api/cert-types"

	// Input JSON
	reqBody := []byte(`{
		"name": "ISO 9001",
		"short_name": "xx"
	}`)

	resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}
		body := string(bodyBytes)

		t.Fatalf("expected 201 Created, got %d with message %s", resp.StatusCode, body)
	}

	var result models.CertificateType

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	// Check UUID format
	if _, err := uuid.Parse(result.Id.String()); err != nil {
		t.Errorf("invalid UUID in response: %v", err)
	}

	// Check CreatedAt / UpdatedAt ≈ now
	now := time.Now().UTC()
	maxDrift := 2 * time.Second

	if result.CreatedAt.Before(now.Add(-maxDrift)) || result.CreatedAt.After(now.Add(maxDrift)) {
		t.Errorf("CreatedAt %v not within expected time window", result.CreatedAt)
	}

	if result.UpdatedAt.Before(now.Add(-maxDrift)) || result.UpdatedAt.After(now.Add(maxDrift)) {
		t.Errorf("UpdatedAt %v not within expected time window", result.UpdatedAt)
	}

	// Check fields match request
	if result.Name != "ISO 9001" {
		t.Errorf("unexpected Name: %s", result.Name)
	}

	StoreCreatedCertificateType(result)
	fmt.Println("Created cert type: " + result.Name + " (" + result.Id.String() + ")")

}

func tHandlerAddCertificate(t *testing.T) {
	// --------------------------------------------------------
	// Arrange
	// --------------------------------------------------------

	url := "http://localhost:8080/api/certificates"

	// Input JSON
	reqBody := []byte(`{
		"cert-type-id": "` + createdCertType.Id.String() + `",
		"cert-number": "ABC123",
		"issuer": "Lloyds",
		"issued-date": "2020-01-01"
	}`)

	// --------------------------------------------------------
	// Act
	// --------------------------------------------------------
	resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	// --------------------------------------------------------
	// Assert – Status
	// --------------------------------------------------------
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}
		body := string(bodyBytes)

		t.Fatalf("expected 201 Created, got %d with message %s", resp.StatusCode, body)
	}

	// --------------------------------------------------------
	// Assert – Body
	// --------------------------------------------------------
	var result models.Certificate

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	// Check UUID format
	if _, err := uuid.Parse(result.ID.String()); err != nil {
		t.Errorf("invalid UUID in response: %v", err)
	}

	// Check CreatedAt / UpdatedAt ≈ now
	now := time.Now().UTC()
	maxDrift := 2 * time.Second

	if result.CreatedAt.Before(now.Add(-maxDrift)) || result.CreatedAt.After(now.Add(maxDrift)) {
		t.Errorf("CreatedAt %v not within expected time window", result.CreatedAt)
	}

	if result.UpdatedAt.Before(now.Add(-maxDrift)) || result.UpdatedAt.After(now.Add(maxDrift)) {
		t.Errorf("UpdatedAt %v not within expected time window", result.UpdatedAt)
	}

	// Check fields match request
	if result.CertType.Name != "ISO 9001" {
		t.Errorf("unexpected Name: %s", result.CertType.Name)
	}
	if result.CertNumber != "ABC123" {
		t.Errorf("unexpected CertNumber: %s", result.CertNumber)
	}
	if result.Issuer != "Lloyds" {
		t.Errorf("unexpected Issuer: %s", result.Issuer)
	}
	StoreCreatedCertificate(result)
}

func tGetCertificateFromId(t *testing.T) {
	if createdCert.ID == uuid.Nil {
		t.Fatalf("POST test did not populate createdCert; ensure StoreCreatedCertificate() is called in the POST test")
	}

	url := fmt.Sprintf("http://localhost:8080/api/certificates/%s", createdCert.ID)

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	var result models.Certificate
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	// Validate exact match
	if result.ID != createdCert.ID {
		t.Errorf("ID mismatch: expected %s got %s", createdCert.ID, result.ID)
	}
	if result.CertType.Name != createdCert.CertType.Name {
		t.Errorf("Name mismatch: expected %s got %s", createdCert.CertType.Name, result.CertType.Name)
	}
	if result.CertNumber != createdCert.CertNumber {
		t.Errorf("CertNumber mismatch: expected %s got %s", createdCert.CertNumber, result.CertNumber)
	}
	if result.Issuer != createdCert.Issuer {
		t.Errorf("Issuer mismatch: expected %s got %s", createdCert.Issuer, result.Issuer)
	}
	if !result.CreatedAt.Equal(createdCert.CreatedAt) {
		t.Errorf("CreatedAt mismatch: expected %v got %v", createdCert.CreatedAt, result.CreatedAt)
	}
	if !result.UpdatedAt.Equal(createdCert.UpdatedAt) {
		t.Errorf("UpdatedAt mismatch: expected %v got %v", createdCert.UpdatedAt, result.UpdatedAt)
	}
	if !result.IssuedDate.Equal(createdCert.IssuedDate) {
		t.Errorf("IssuedDate mismatch")
	}

}

func tGetAllCerts(t *testing.T) {

	expectedCount := 6

	resp, errResp := http.Get("http://localhost:8080/api/certificates")
	if errResp != nil {
		t.Fatalf("failed to call /api/certificates: %v", errResp)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// ---- Unmarshal returned certificates ----
	var returned []models.Certificate
	if err := json.Unmarshal(body, &returned); err != nil {
		t.Fatalf("failed to unmarshal API response: %v", err)
	}

	// ---- List returned certificates to stdout ----
	fmt.Println("=== Certificates returned by API ===")
	for i, c := range returned {
		fmt.Printf("%d: %s (%s) issued by %s on %s\n",
			i, c.CertType.Name, c.CertNumber, c.Issuer, c.IssuedDate)
	}
	fmt.Println("====================================")

	// ---- Check count ----
	if len(returned) != expectedCount {
		t.Fatalf("certificate count mismatch: got %d, expected %d",
			len(returned), expectedCount)
	}

}

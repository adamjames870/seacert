package certificates

import (
	"database/sql"
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/google/uuid"
)

func TestMapCertificateViewDbToDomain(t *testing.T) {
	id := uuid.New()
	certTypeId := uuid.New()
	issuerId := uuid.New()
	now := time.Now().UTC()

	view := sqlc.CertView{
		ID:                    id,
		CreatedAt:             now,
		UpdatedAt:             now,
		CertNumber:            "CERT-123",
		IssuedDate:            now,
		AlternativeName:       sql.NullString{String: "Alt Name", Valid: true},
		Remarks:               sql.NullString{String: "Remarks", Valid: true},
		ManualExpiry:          sql.NullTime{Time: now.AddDate(1, 0, 0), Valid: true},
		DocumentPath:          sql.NullString{String: "path/to/doc", Valid: true},
		CertTypeID:            certTypeId,
		CertTypeCreatedAt:     now,
		CertTypeUpdatedAt:     now,
		CertTypeName:          "Cert Type Name",
		CertTypeShortName:     "CTN",
		CertTypeStcwReference: sql.NullString{String: "VI/1", Valid: true},
		NormalValidityMonths:  sql.NullInt32{Int32: 60, Valid: true},
		IssuerID:              issuerId,
		IssuerCreatedAt:       now,
		IssuerUpdatedAt:       now,
		IssuerName:            "Issuer Name",
		IssuerCountry:         sql.NullString{String: "UK", Valid: true},
		IssuerWebsite:         sql.NullString{String: "https://mca.gov.uk", Valid: true},
		Deleted:               false,
		HasSuccessor:          true,
		HasPredecessors:       false,
	}

	got := MapCertificateViewDbToDomain(view)

	if got.Id != id {
		t.Errorf("expected ID %v, got %v", id, got.Id)
	}
	if got.CertNumber != "CERT-123" {
		t.Errorf("expected CertNumber CERT-123, got %s", got.CertNumber)
	}
	if got.CertType.Id != certTypeId {
		t.Errorf("expected CertTypeId %v, got %v", certTypeId, got.CertType.Id)
	}
	if got.Issuer.Id != issuerId {
		t.Errorf("expected IssuerId %v, got %v", issuerId, got.Issuer.Id)
	}
	if got.HasSuccessors != view.HasSuccessor {
		t.Errorf("expected HasSuccessors %v, got %v", view.HasSuccessor, got.HasSuccessors)
	}
	if got.DocumentPath != "path/to/doc" {
		t.Errorf("expected DocumentPath path/to/doc, got %s", got.DocumentPath)
	}
}

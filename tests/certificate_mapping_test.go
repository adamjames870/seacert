package tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/domain/certificates"
	"github.com/adamjames870/seacert/internal/domain/issuers"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func TestCertificateMapping(t *testing.T) {
	certId := uuid.New()
	certTypeId := uuid.New()
	issuerId := uuid.New()
	now := time.Now()

	certType := cert_types.CertificateType{
		Id:        certTypeId,
		Name:      "Test Cert Type",
		ShortName: "TCT",
	}

	issuer := issuers.Issuer{
		Id:      issuerId,
		Name:    "Test Issuer",
		Country: "Test Country",
	}

	dbCert := sqlc.Certificate{
		ID:              certId,
		CreatedAt:       now,
		UpdatedAt:       now,
		CertTypeID:      certTypeId,
		CertNumber:      "12345",
		IssuerID:        issuerId,
		IssuedDate:      now,
		AlternativeName: sql.NullString{String: "Alt Name", Valid: true},
		Remarks:         sql.NullString{String: "Some remarks", Valid: true},
	}

	t.Run("MapCertificateDbToDomain", func(t *testing.T) {
		got := certificates.MapCertificateDbToDomain(dbCert, certType, issuer)
		if got.Id != certId {
			t.Errorf("expected ID %v, got %v", certId, got.Id)
		}
		if got.CertNumber != "12345" {
			t.Errorf("expected CertNumber 12345, got %s", got.CertNumber)
		}
		if got.AlternativeName != "Alt Name" {
			t.Errorf("expected AlternativeName Alt Name, got %s", got.AlternativeName)
		}
	})

	t.Run("MapCertificateDomainToDto", func(t *testing.T) {
		domainCert := certificates.Certificate{
			Id:         certId,
			CertNumber: "12345",
			CertType:   certType,
			Issuer:     issuer,
		}
		got := certificates.MapCertificateDomainToDto(domainCert)
		if got.Id != certId.String() {
			t.Errorf("expected ID %s, got %s", certId.String(), got.Id)
		}
		if got.CertTypeName != certType.Name {
			t.Errorf("expected CertTypeName %s, got %s", certType.Name, got.CertTypeName)
		}
	})

	t.Run("MapCertificateDtoToDomain", func(t *testing.T) {
		dtoCert := dto.Certificate{
			Id:           certId.String(),
			CertNumber:   "12345",
			CertTypeId:   certTypeId.String(),
			CertTypeName: certType.Name,
			IssuerName:   issuer.Name,
		}
		got := certificates.MapCertificateDtoToDomain(dtoCert)
		if got.Id != certId {
			t.Errorf("expected ID %v, got %v", certId, got.Id)
		}
		if got.CertNumber != "12345" {
			t.Errorf("expected CertNumber 12345, got %s", got.CertNumber)
		}
	})

	t.Run("MapCertificateDomainToDb", func(t *testing.T) {
		domainCert := certificates.Certificate{
			Id:         certId,
			CertNumber: "12345",
			CertType:   certType,
			Issuer:     issuer,
		}
		got := certificates.MapCertificateDomainToDb(domainCert)
		if got.ID != certId {
			t.Errorf("expected ID %v, got %v", certId, got.ID)
		}
		if got.CertNumber != "12345" {
			t.Errorf("expected CertNumber 12345, got %s", got.CertNumber)
		}
	})
}

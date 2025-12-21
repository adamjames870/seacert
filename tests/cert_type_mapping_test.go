package tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/google/uuid"
)

func TestCertTypeMapping(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	dbCertType := sqlc.CertificateType{
		ID:                   id,
		CreatedAt:            now,
		UpdatedAt:            now,
		Name:                 "Master Mariner",
		ShortName:            "MM",
		StcwReference:        sql.NullString{String: "II/2", Valid: true},
		NormalValidityMonths: sql.NullInt32{Int32: 60, Valid: true},
	}

	// DB to Domain
	domainCertType := cert_types.MapCertificateTypeDbToDomain(dbCertType)
	if domainCertType.Id != id || domainCertType.Name != "Master Mariner" || domainCertType.StcwReference != "II/2" {
		t.Errorf("DB to Domain mapping failed: %+v", domainCertType)
	}

	// Domain to DTO
	dtoCertType := cert_types.MapCertificateTypeDomainToDto(domainCertType)
	if dtoCertType.Id != id.String() || dtoCertType.Name != "Master Mariner" || dtoCertType.StcwRef != "II/2" {
		t.Errorf("Domain to DTO mapping failed: %+v", dtoCertType)
	}

	// DTO to Domain
	domainCertType2 := cert_types.MapCertificateTypeDtoToDomain(dtoCertType)
	if domainCertType2.Id != id || domainCertType2.Name != "Master Mariner" || domainCertType2.StcwReference != "II/2" {
		t.Errorf("DTO to Domain mapping failed: %+v", domainCertType2)
	}

	// Domain to DB
	dbCertType2 := cert_types.MapCertificateTypeDomainToDb(domainCertType2)
	if dbCertType2.ID != id || dbCertType2.Name != "Master Mariner" || dbCertType2.StcwReference.String != "II/2" {
		t.Errorf("Domain to DB mapping failed: %+v", dbCertType2)
	}
}

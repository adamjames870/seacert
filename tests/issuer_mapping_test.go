package tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain/issuers"
	"github.com/google/uuid"
)

func TestIssuerMapping(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	dbIssuer := sqlc.Issuer{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "MCA",
		Country:   sql.NullString{String: "UK", Valid: true},
		Website:   sql.NullString{String: "https://mca.gov.uk", Valid: true},
	}

	// DB to Domain
	domainIssuer := issuers.MapIssuerDbToDomain(dbIssuer)
	if domainIssuer.Id != id || domainIssuer.Name != "MCA" || domainIssuer.Country != "UK" {
		t.Errorf("DB to Domain mapping failed: %+v", domainIssuer)
	}

	// Domain to DTO
	dtoIssuer := issuers.MapIssuerDomainToDto(domainIssuer)
	if dtoIssuer.Id != id.String() || dtoIssuer.Name != "MCA" || dtoIssuer.Country != "UK" {
		t.Errorf("Domain to DTO mapping failed: %+v", dtoIssuer)
	}

	// DTO to Domain
	domainIssuer2 := issuers.MapIssuerDtoToDomain(dtoIssuer)
	if domainIssuer2.Id != id || domainIssuer2.Name != "MCA" || domainIssuer2.Country != "UK" {
		t.Errorf("DTO to Domain mapping failed: %+v", domainIssuer2)
	}

	// Domain to DB
	dbIssuer2 := issuers.MapIssuerDomainToDb(domainIssuer2)
	if dbIssuer2.ID != id || dbIssuer2.Name != "MCA" || dbIssuer2.Country.String != "UK" {
		t.Errorf("Domain to DB mapping failed: %+v", dbIssuer2)
	}
}

package issuers

import (
	"database/sql"
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func TestMapIssuer(t *testing.T) {
	id := uuid.New()
	now := time.Now().UTC()

	dbIssuer := sqlc.Issuer{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "MCA",
		Country: sql.NullString{
			String: "UK",
			Valid:  true,
		},
		Website: sql.NullString{
			String: "https://www.gov.uk/mca",
			Valid:  true,
		},
	}

	domainIssuer := Issuer{
		Id:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "MCA",
		Country:   "UK",
		Website:   "https://www.gov.uk/mca",
	}

	dtoIssuer := dto.Issuer{
		Id:        id.String(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "MCA",
		Country:   "UK",
		Website:   "https://www.gov.uk/mca",
	}

	t.Run("MapIssuerDbToDomain", func(t *testing.T) {
		got := MapIssuerDbToDomain(dbIssuer)
		if got != domainIssuer {
			t.Errorf("expected %+v, got %+v", domainIssuer, got)
		}
	})

	t.Run("MapIssuerDomainToDto", func(t *testing.T) {
		got := MapIssuerDomainToDto(domainIssuer)
		if got != dtoIssuer {
			t.Errorf("expected %+v, got %+v", dtoIssuer, got)
		}
	})

	t.Run("MapIssuerDtoToDomain", func(t *testing.T) {
		got := MapIssuerDtoToDomain(dtoIssuer)
		if got != domainIssuer {
			t.Errorf("expected %+v, got %+v", domainIssuer, got)
		}
	})

	t.Run("MapIssuerDomainToDb", func(t *testing.T) {
		got := MapIssuerDomainToDb(domainIssuer)
		if got.ID != dbIssuer.ID ||
			got.Name != dbIssuer.Name ||
			got.Country.String != dbIssuer.Country.String ||
			got.Website.String != dbIssuer.Website.String {
			t.Errorf("expected %+v, got %+v", dbIssuer, got)
		}
	})
}
